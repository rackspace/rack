package auth

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/rackspace"
)

// reauthFunc is what the ServiceClient uses to re-authenticate.
func reauthFunc(pc *gophercloud.ProviderClient, ao gophercloud.AuthOptions) func() error {
	return func() error {
		return rackspace.AuthenticateV2(pc, ao)
	}
}

// NewClient creates and returns a Rackspace client for the given service.
func NewClient(c *cli.Context, serviceType string) (*gophercloud.ServiceClient, error) {
	// get the user's authentication credentials
	ao, region, err := Credentials(c)
	if err != nil {
		return nil, err
	}

	if c.GlobalIsSet("no-cache") || c.IsSet("no-cache") {
		return authFromScratch(*ao, region, serviceType)
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
		pc, err := rackspace.NewClient(ao.IdentityEndpoint)
		if err == nil {
			pc.TokenID = creds.TokenID
			pc.ReauthFunc = reauthFunc(pc, *ao)
			pc.UserAgent.Prepend(util.UserAgent)
			return &gophercloud.ServiceClient{
				ProviderClient: pc,
				Endpoint:       creds.ServiceEndpoint,
			}, nil
		}
	} else {
		return authFromScratch(*ao, region, serviceType)
	}

	return nil, nil
}

func authFromScratch(ao gophercloud.AuthOptions, region, serviceType string) (*gophercloud.ServiceClient, error) {
	pc, err := rackspace.AuthenticatedClient(ao)
	if err != nil {
		return nil, err
	}
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
		return nil, fmt.Errorf("Unable to create service client: Unknown service type: %s", serviceType)
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
func Credentials(c *cli.Context) (*gophercloud.AuthOptions, string, error) {
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

	if len(need) > 0 {
		haveString := ""
		for k, v := range have {
			haveString += fmt.Sprintf("%s: %s (from %s)\n", k, v.value, v.from)
		}

		needString := ""
		for k := range need {
			needString += fmt.Sprintf("%s\n", k)
		}

		authErrSlice := []string{"There are some required Rackspace Cloud credentials that we couldn't find.",
			"Here's what we have:",
			fmt.Sprintf("%s", haveString),
			"and here's what we we're missing:",
			fmt.Sprintf("%s", needString),
			"",
			"You can set any of these credentials in the following ways:",
			"- Run `rack config` to interactively create a configuration file,",
			"- Specify it in the command as a flag (--username, --apikey, --region), or",
			"- Export it as an environment variable (RS_USERNAME, RS_API_KEY, RS_REGION_NAME).",
			"",
		}

		return nil, "", fmt.Errorf(strings.Join(authErrSlice, "\n"))
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
