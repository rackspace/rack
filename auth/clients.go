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
func NewClient(c *cli.Context, t string) (*gophercloud.ServiceClient, error) {
	// get the user's authentication credentials
	ao, region := Credentials(c)
	// upper-case the region
	region = strings.ToUpper(region)
	// allow Gophercloud to re-authenticate
	ao.AllowReauth = true
	// if the user didn't provide an auth URL, default to the Rackspace US endpoint
	if ao.IdentityEndpoint == "" {
		ao.IdentityEndpoint = rackspace.RackspaceUSIdentity
	}

	var sc *gophercloud.ServiceClient

	// form the cache key
	cacheKey := CacheKey(ao, region, t)
	// initialize cache
	cache := &Cache{}
	// get the value from the cache
	creds, err := cache.Value(cacheKey)
	// if there was an error accessing the cache or there was nothing in the cache,
	// authenticate from scratch
	if err != nil || creds == nil {
		pc, err := rackspace.AuthenticatedClient(ao)
		if err != nil {
			return nil, fmt.Errorf("Error creating ProviderClient: %s\n", err)
		}
		switch t {
		case "compute":
			sc, err = rackspace.NewComputeV2(pc, gophercloud.EndpointOpts{
				Region: region,
			})
			break
		case "blockstorage":
			sc, err = rackspace.NewBlockStorageV1(pc, gophercloud.EndpointOpts{
				Region: region,
			})
			break
		case "networking":
			sc, err = rackspace.NewNetworkV2(pc, gophercloud.EndpointOpts{
				Region: region,
			})
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error creating ServiceClient: %s\n", err)
		}
	} else {
		// we successfully retrieved a value from the cache
		pc := &gophercloud.ProviderClient{
			IdentityBase:     creds.IdentityBase,
			IdentityEndpoint: creds.IdentityEndpoint,
			TokenID:          creds.TokenID,
			HTTPClient:       creds.HTTPClient,
		}
		pc.ReauthFunc = reauthFunc(pc, ao)
		sc = &gophercloud.ServiceClient{
			ProviderClient: pc,
			Endpoint:       creds.ServiceEndpoint,
		}
	}

	// set the user-agent
	sc.UserAgent.Prepend("rack/" + util.Version)

	return sc, nil
}

// Credentials determines the appropriate authentication method for the user.
// It returns a gophercloud.AuthOptions object and a region.
//
// It will use command-line authentication parameters if available, then it will
// look for any unset parameters in the config file, and then finally in
// environment variables.
func Credentials(c *cli.Context) (gophercloud.AuthOptions, string) {
	have := make(map[string]string)
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
		configfile(c, have, need)
		// still unset auth variables?
		if len(need) != 0 {
			// if so, look in environment variables
			envvars(have, need)
		}
	}

	ao := gophercloud.AuthOptions{
		Username:         have["username"],
		APIKey:           have["apikey"],
		IdentityEndpoint: have["authurl"],
	}
	return ao, have["region"]
}
