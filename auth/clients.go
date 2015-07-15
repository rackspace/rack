package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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
func NewClient(c *cli.Context, serviceType string, logger *logrus.Logger) (*gophercloud.ServiceClient, error) {
	// get the user's authentication credentials
	ao, region, err := Credentials(c, logger)
	if err != nil {
		return nil, err
	}

	if c.GlobalIsSet("no-cache") || c.IsSet("no-cache") {
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

type authCred struct {
	value string
	from  string
}

// Credentials determines the appropriate authentication method for the user.
// It returns a gophercloud.AuthOptions object and a region.
//
// It will use command-line authentication parameters if available, then it will
// look for any unset parameters in the config file, and then finally in
// environment variables.
func Credentials(c *cli.Context, logger *logrus.Logger) (*gophercloud.AuthOptions, string, error) {
	have := make(map[string]authCred)
	need := map[string]string{
		"username": "",
		"apikey":   "",
		"authurl":  "",
		"region":   "",
	}

	// use command-line options if available
	cliopts(c, have, need)
	// are there any unset auth variables?
	if len(need) != 0 {
		// if so, look in config file
		err := configfile(c, have, need)
		if err != nil {
			return nil, "", err
		}
		// still unset auth variables?
		if len(need) != 0 {
			// if so, look in environment variables
			envvars(have, need)
		}
	}

	// if the user didn't provide an auth URL, default to the Rackspace US endpoint
	if _, ok := have["authurl"]; !ok || have["authurl"].value == "" {
		have["authurl"] = authCred{value: rackspace.RackspaceUSIdentity, from: "default value"}
		delete(need, "authurl")
	}

	haveString := ""
	for k, v := range have {
		haveString += fmt.Sprintf("%s: %s (from %s)\n", k, v.value, v.from)
	}

	if len(need) > 0 {
		needString := ""
		for k := range need {
			needString += fmt.Sprintf("%s\n", k)
		}

		authErrSlice := []string{"There are some required Rackspace Cloud credentials that we couldn't find.",
			"Here's what we have:",
			fmt.Sprintf("%s", haveString),
			"and here's what we're missing:",
			fmt.Sprintf("%s", needString),
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

	ao := &gophercloud.AuthOptions{
		Username:         have["username"].value,
		APIKey:           have["apikey"].value,
		IdentityEndpoint: have["authurl"].value,
	}

	// upper-case the region
	region := strings.ToUpper(have["region"].value)
	// allow Gophercloud to re-authenticate
	ao.AllowReauth = true

	return ao, region, nil
}

// LogRoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default Gophercloud RoundTripper to allow for logging.
type LogRoundTripper struct {
	Logger *logrus.Logger
	rt     http.RoundTripper
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
