package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/internal/github.com/Sirupsen/logrus"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace"
	"github.com/jrperritt/rack/util"
)

// reauthFunc is what the ServiceClient uses to re-authenticate.
func reauthFunc(pc *gophercloud.ProviderClient, ao gophercloud.AuthOptions) func() error {
	return func() error {
		return rackspace.AuthenticateV2(pc, ao)
	}
}

// NewClient creates and returns a Rackspace client for the given service.
func NewClient(c *cli.Context, serviceType string, logger *logrus.Logger, noCache bool) (*gophercloud.ServiceClient, error) {
	// get the user's authentication credentials
	ao, region, err := Credentials(c, logger)
	if err != nil {
		return nil, err
	}

	if noCache {
		return authFromScratch(*ao, region, serviceType, logger)
	}

	// form the cache key
	cacheKey := CacheKey(*ao, region, serviceType)
	// initialize cache
	cache := &Cache{}
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
		return authFromScratch(*ao, region, serviceType, logger)
	}

	return nil, nil
}

func authFromScratch(ao gophercloud.AuthOptions, region, serviceType string, logger *logrus.Logger) (*gophercloud.ServiceClient, error) {
	logger.Info("Not using cache; Authenticating from scratch.\n")
	pc, err := rackspace.AuthenticatedClient(ao)
	if err != nil {
		return nil, err
	}
	pc.HTTPClient = newHTTPClient()
	var sc *gophercloud.ServiceClient
	switch serviceType {
	case "compute":
		sc, err = rackspace.NewComputeV2(pc, gophercloud.EndpointOpts{
			Region: region,
		})
		break
	case "object-store":
		sc, err = rackspace.NewObjectStorageV1(pc, gophercloud.EndpointOpts{
			Region: region,
		})
		break
	case "blockstorage":
		sc, err = rackspace.NewBlockStorageV1(pc, gophercloud.EndpointOpts{
			Region: region,
		})
		break
	case "network":
		sc, err = rackspace.NewNetworkV2(pc, gophercloud.EndpointOpts{
			Region: region,
		})
		break
	}
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, fmt.Errorf("Unable to create service client: Unknown service type: %s\n", serviceType)
	}
	sc.UserAgent.Prepend(util.UserAgent)
	return sc, nil
}

// Credentials determines the appropriate authentication method for the user.
// It returns a gophercloud.AuthOptions object and a region.
//
// It will use command-line authentication parameters if available, then it will
// look for any unset parameters in the config file, and then finally in
// environment variables.
func Credentials(c *cli.Context, logger *logrus.Logger) (*gophercloud.AuthOptions, string, error) {
	have := make(map[string]commandoptions.Cred)
	want := map[string]string{
		"username":   "",
		"apikey":     "",
		"authurl":    "",
		"region":     "",
		"tenant-id":  "",
		"auth-token": "",
	}

	// use command-line options if available
	commandoptions.CLIopts(c, have, want)
	// are there any unset auth variables?
	if len(want) != 0 {
		// if so, look in config file
		err := commandoptions.ConfigFile(c, have, want)
		if err != nil {
			return nil, "", err
		}
		// still unset auth variables?
		if len(want) != 0 {
			// if so, look in environment variables
			envvars(have, want)
		}
	}

	var ao *gophercloud.AuthOptions

	if _, ok := have["username"]; ok {
		if _, ok := have["apikey"]; ok {
			delete(want, "tenant-id")
			delete(want, "auth-token")
			ao = &gophercloud.AuthOptions{
				Username: have["username"].Value,
				APIKey:   have["apikey"].Value,
			}
		} else {
			return nil, "", fmt.Errorf("You must provide the `apikey` flag with the `username` flag.")
		}
	} else if _, ok := have["tenant-id"]; ok {
		if _, ok := have["auth-token"]; ok {
			delete(want, "username")
			delete(want, "api-key")
			ao = &gophercloud.AuthOptions{
				TenantID: have["tenant-id"].Value,
				Token:    have["auth-token"].Value,
			}
		} else {
			return nil, "", fmt.Errorf("You must provide the `auth-token` flag with the `tenant-id` flag.")
		}
	}

	// if the user didn't provide an auth URL, default to the Rackspace US endpoint
	if _, ok := have["authurl"]; !ok || have["authurl"].Value == "" {
		have["authurl"] = commandoptions.Cred{Value: rackspace.RackspaceUSIdentity, From: "default value"}
		delete(want, "authurl")
	}

	haveString := ""
	for k, v := range have {
		haveString += fmt.Sprintf("%s: %s (from %s)\n", k, v.Value, v.From)
	}

	if len(want) > 0 {
		wantString := ""
		for k := range want {
			wantString += fmt.Sprintf("%s\n", k)
		}

		authErrSlice := []string{"There are some required Rackspace Cloud credentials that we couldn't find.",
			"Here's what we have:",
			fmt.Sprintf("%s", haveString),
			"and here's what we're missing:",
			fmt.Sprintf("%s", wantString),
			"",
			"You can set any of these credentials in the following ways:",
			"- Run `rack configure` to interactively create a configuration file,",
			"- Specify it in the command as a flag (--username, --apikey, --region), or",
			"- Export it as an environment variable (RS_USERNAME, RS_API_KEY, RS_REGION_NAME).",
			"",
		}

		return nil, "", fmt.Errorf(strings.Join(authErrSlice, "\n"))
	}

	if logger != nil {
		logger.Infof("Authentication Credentials:\n%s\n", haveString)
	}

	// upper-case the region
	region := strings.ToUpper(have["region"].Value)
	// allow Gophercloud to re-authenticate
	ao.AllowReauth = true
	ao.IdentityEndpoint = have["authurl"].Value

	return ao, region, nil
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
		request.Body, err = lrt.logRequestBody(request.Body, request.Header)
		if err != nil {
			return nil, err
		}
	}

	lrt.Logger.Infof("Request URL: %s\n", request.URL)

	response, err := lrt.rt.RoundTrip(request)
	if response.StatusCode == http.StatusUnauthorized {
		if lrt.numReauthAttempts == 3 {
			return response, fmt.Errorf("Tried to re-authenticate 3 times with no success.")
		}
		lrt.numReauthAttempts++
	}
	if err != nil {
		return response, err
	}

	lrt.Logger.Debugf("Response Status: %s\n", response.Status)

	info, err := json.MarshalIndent(response.Header, "", "  ")
	if err != nil {
		lrt.Logger.Debugf(fmt.Sprintf("Error logging request: %s\n", err))
	}
	lrt.Logger.Debugf("Response Headers: %+v\n", string(info))

	return response, nil
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
