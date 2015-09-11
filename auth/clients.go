package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/internal/github.com/Sirupsen/logrus"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	tokens2 "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace"
	"github.com/rackspace/rack/util"
)

var usernameAuthErrSlice = []string{"There are some required Rackspace Cloud credentials that we couldn't find.",
	"Here's what we have:",
	"%s",
	"and here's what we're missing:",
	"%s",
	"",
	"You can set any of these credentials in the following ways:",
	"- Run `rack configure` to interactively create a configuration file,",
	"- Specify it in the command as a flag (--username, --api-key), or",
	"- Export it as an environment variable (RS_USERNAME, RS_API_KEY).",
	"",
}

var tenantIDAuthErrSlice = []string{"There are some required Rackspace Cloud credentials that we couldn't find.",
	"Here's what we have:",
	"%s",
	"and here's what we're missing:",
	"%s",
	"",
	"You can set the missing credentials with command-line flags (--auth-token, --auth-tenant-id)",
	"",
}

// Err returns the custom error to print when authentication fails.
func Err(have map[string]commandoptions.Cred, want map[string]string, errMsg []string) error {
	haveString := ""
	for k, v := range have {
		haveString += fmt.Sprintf("%s: %s (from %s)\n", k, v.Value, v.From)
	}

	if len(want) > 0 {
		wantString := ""
		for k := range want {
			wantString += fmt.Sprintf("%s\n", k)
		}

		return fmt.Errorf(fmt.Sprintf(strings.Join(errMsg, "\n"), haveString, wantString))
	}

	return nil
}

// CredentialsResult holds the information acquired from looking for authentication
// credentials.
type CredentialsResult struct {
	AuthOpts *gophercloud.AuthOptions
	Region   string
	Have     map[string]commandoptions.Cred
	Want     map[string]string
}

func findAuthOpts(c *cli.Context, have map[string]commandoptions.Cred, want map[string]string) error {
	// use command-line options if available
	commandoptions.CLIopts(c, have, want)
	// are there any unset auth variables?
	if len(want) != 0 {
		// if so, look in config file
		err := commandoptions.ConfigFile(c, have, want)
		if err != nil {
			return err
		}
		// still unset auth variables?
		if len(want) != 0 {
			// if so, look in environment variables
			envvars(have, want)
		}
	}

	return nil
}

// reauthFunc is what the ServiceClient uses to re-authenticate.
func reauthFunc(pc *gophercloud.ProviderClient, ao gophercloud.AuthOptions) func() error {
	return func() error {
		return rackspace.AuthenticateV2(pc, ao)
	}
}

// NewClient creates and returns a Rackspace client for the given service.
func NewClient(c *cli.Context, serviceType string, logger *logrus.Logger, noCache bool, useServiceNet bool) (*gophercloud.ServiceClient, error) {
	// get the user's authentication credentials
	credsResult, err := Credentials(c, logger)
	if err != nil {
		return nil, err
	}

	logMsg := "Using public endpoint"
	urlType := gophercloud.AvailabilityPublic
	if useServiceNet {
		logMsg = "Using service net endpoint"
		urlType = gophercloud.AvailabilityInternal
	}
	logger.Infoln(logMsg)

	if noCache {
		return authFromScratch(credsResult, serviceType, urlType, logger)
	}

	ao := credsResult.AuthOpts
	region := credsResult.Region

	// form the cache key
	cacheKey := CacheKey(*ao, region, serviceType, urlType)
	// initialize cache
	cache := &Cache{}
	logger.Infof("Looking in the cache for cache key: %s\n", cacheKey)
	// get the value from the cache
	creds, err := cache.Value(cacheKey)
	// if there was an error accessing the cache or there was nothing in the cache,
	// authenticate from scratch
	if err == nil && creds != nil {
		// we successfully retrieved a value from the cache
		logger.Infof("Using token from cache: %s\n", creds.TokenID)
		pc, err := rackspace.NewClient(ao.IdentityEndpoint)
		if err == nil {
			pc.TokenID = creds.TokenID
			pc.ReauthFunc = reauthFunc(pc, *ao)
			pc.UserAgent.Prepend(util.UserAgent)
			pc.HTTPClient = newHTTPClient()
			return &gophercloud.ServiceClient{
				ProviderClient: pc,
				Endpoint:       creds.ServiceEndpoint,
			}, nil
		}
	} else {
		return authFromScratch(credsResult, serviceType, urlType, logger)
	}

	return nil, nil
}

func authFromScratch(credsResult *CredentialsResult, serviceType string, urlType gophercloud.Availability, logger *logrus.Logger) (*gophercloud.ServiceClient, error) {
	logger.Info("Not using cache; Authenticating from scratch.\n")

	ao := credsResult.AuthOpts
	region := credsResult.Region

	pc, err := rackspace.AuthenticatedClient(*ao)
	if err != nil {
		switch err.(type) {
		case *tokens2.ErrNoPassword:
			return nil, errors.New("Please supply an API key.")
		}
		return nil, err
	}
	pc.HTTPClient = newHTTPClient()
	var sc *gophercloud.ServiceClient
	switch serviceType {
	case "compute":
		sc, err = rackspace.NewComputeV2(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	case "object-store":
		sc, err = rackspace.NewObjectStorageV1(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	case "blockstorage":
		sc, err = rackspace.NewBlockStorageV1(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	case "network":
		sc, err = rackspace.NewNetworkV2(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	case "orchestration":
		sc, err = rackspace.NewOrchestrationV1(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	}
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, fmt.Errorf("Unable to create service client: Unknown service type: %s\n", serviceType)
	}
	if sc.Endpoint == "/" {
		return nil, fmt.Errorf(strings.Join([]string{"You wanted to use service net for the %s request",
			"but the %s service doesn't have an internal URL.\n"}, " "), serviceType, serviceType)
	}
	logger.Debugf("Created %s service client: %+v", serviceType, sc)
	sc.UserAgent.Prepend(util.UserAgent)
	return sc, nil
}

// Credentials determines the appropriate authentication method for the user.
// It returns a gophercloud.AuthOptions object and a region.
//
// It will use command-line authentication parameters if available, then it will
// look for any unset parameters in the config file, and then finally in
// environment variables.
func Credentials(c *cli.Context, logger *logrus.Logger) (*CredentialsResult, error) {
	ao := &gophercloud.AuthOptions{
		AllowReauth: true,
	}

	have := make(map[string]commandoptions.Cred)

	// let's looks for a region and identity endpoint
	want := map[string]string{
		"auth-url": "",
		"region":   "",
	}

	err := findAuthOpts(c, have, want)
	if err != nil {
		return nil, err
	}

	// if the user didn't provide an auth URL, default to the Rackspace US endpoint
	if _, ok := have["auth-url"]; !ok || have["auth-url"].Value == "" {
		have["auth-url"] = commandoptions.Cred{Value: rackspace.RackspaceUSIdentity, From: "default value"}
		delete(want, "auth-url")
	}
	ao.IdentityEndpoint = have["auth-url"].Value

	// upper-case the region
	region := strings.ToUpper(have["region"].Value)
	delete(want, "region")

	// now we check for token authentication (only allowed via the command-line)
	want["auth-tenant-id"] = ""
	want["auth-token"] = ""
	commandoptions.CLIopts(c, have, want)

	// if a tenant ID was provided on the command-line, we don't bother checking for a
	// username or api key
	if have["auth-tenant-id"].Value != "" || have["auth-token"].Value != "" {
		if tenantID, ok := have["auth-tenant-id"]; ok {
			ao.TenantID = tenantID.Value
			ao.TokenID = have["auth-token"].Value
			delete(want, "auth-token")
		} else {
			return nil, Err(have, want, tenantIDAuthErrSlice)
		}
	} else {
		// otherwise, let's look for a username and API key
		want = map[string]string{
			"username": "",
			"api-key":  "",
		}
		err = findAuthOpts(c, have, want)
		if err != nil {
			return nil, err
		}
		if have["username"].Value != "" || have["api-key"].Value != "" {
			if username, ok := have["username"]; ok {
				ao.Username = username.Value
				ao.APIKey = have["api-key"].Value
				delete(want, "api-key")
			} else {
				return nil, Err(have, want, usernameAuthErrSlice)
			}
		} else {
			return nil, Err(have, want, usernameAuthErrSlice)
		}
	}

	if logger != nil {
		haveString := ""
		for k, v := range have {
			haveString += fmt.Sprintf("%s: %s (from %s)\n", k, v.Value, v.From)
		}
		logger.Infof("Authentication Credentials:\n%s\n", haveString)
	}

	credsResult := &CredentialsResult{
		AuthOpts: ao,
		Region:   region,
		Have:     have,
		Want:     want,
	}

	return credsResult, nil
}

// LogRoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default Gophercloud RoundTripper to allow for logging.
type LogRoundTripper struct {
	Logger            *logrus.Logger
	rt                http.RoundTripper
	numReauthAttempts int
}

// newHTTPClient return a custom HTTP client that allows for logging relevant
// information before and after the HTTP request.
func newHTTPClient() http.Client {
	return http.Client{
		Transport: &LogRoundTripper{
			rt: http.DefaultTransport,
		},
	}
}

// RoundTrip performs a round-trip HTTP request and logs relevant information about it.
func (lrt *LogRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var err error

	if lrt.Logger.Level == logrus.DebugLevel && request.Body != nil {
		fmt.Println("Logging request body...")
		request.Body, err = lrt.logRequestBody(request.Body, request.Header)
		if err != nil {
			return nil, err
		}
	}

	info, err := json.MarshalIndent(request.Header, "", "  ")
	if err != nil {
		lrt.Logger.Debugf(fmt.Sprintf("Error logging request headers: %s\n", err))
	}
	lrt.Logger.Debugf("Request Headers: %+v\n", string(info))

	lrt.Logger.Infof("Request URL: %s\n", request.URL)

	response, err := lrt.rt.RoundTrip(request)
	if response == nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		if lrt.numReauthAttempts == 3 {
			return response, fmt.Errorf("Tried to re-authenticate 3 times with no success.")
		}
		lrt.numReauthAttempts++
	}

	lrt.Logger.Debugf("Response Status: %s\n", response.Status)

	info, err = json.MarshalIndent(response.Header, "", "  ")
	if err != nil {
		lrt.Logger.Debugf(fmt.Sprintf("Error logging response headers: %s\n", err))
	}
	lrt.Logger.Debugf("Response Headers: %+v\n", string(info))

	return response, err
}

func (lrt *LogRoundTripper) logRequestBody(original io.ReadCloser, headers http.Header) (io.ReadCloser, error) {
	defer original.Close()

	var bs bytes.Buffer
	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, err
	}

	contentType := headers.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		debugInfo := lrt.formatJSON(bs.Bytes())
		lrt.Logger.Debugf("Request Options: %s\n", debugInfo)
	} else {
		lrt.Logger.Debugf("Request Options: %s\n", bs.String())
	}

	return ioutil.NopCloser(strings.NewReader(bs.String())), nil
}

func (lrt *LogRoundTripper) formatJSON(raw []byte) string {
	var data map[string]interface{}

	err := json.Unmarshal(raw, &data)
	if err != nil {
		lrt.Logger.Debugf("Unable to parse JSON: %s\n\n", err)
		return string(raw)
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		lrt.Logger.Debugf("Unable to re-marshal JSON: %s\n", err)
		return string(raw)
	}

	return string(pretty)
}
