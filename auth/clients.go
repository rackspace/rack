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
	// upper-case the region
	region = strings.ToUpper(region)
	// allow Gophercloud to re-authenticate
	ao.AllowReauth = true
	// if the user didn't provide an auth URL, default to the Rackspace US endpoint
	if ao.IdentityEndpoint == "" {
		ao.IdentityEndpoint = rackspace.RackspaceUSIdentity
	}

	if c.GlobalIsSet("no-cache") {
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
			pc.UserAgent.Prepend("rack/" + util.Version)
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
		return nil, fmt.Errorf("Error creating ProviderClient: %s\n", err)
	}
	var sc *gophercloud.ServiceClient
	switch serviceType {
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
	if sc == nil {
		return nil, fmt.Errorf("Unable to create service client: Unknown service type: %s", serviceType)
	}
	sc.UserAgent.Prepend("rack/" + util.Version)
	return sc, nil
}

// Credentials determines the appropriate authentication method for the user.
// It returns a gophercloud.AuthOptions object and a region.
//
// It will use command-line authentication parameters if available, then it will
// look for any unset parameters in the config file, and then finally in
// environment variables.
func Credentials(c *cli.Context) (*gophercloud.AuthOptions, string, error) {
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

	ao := &gophercloud.AuthOptions{
		Username:         have["username"],
		APIKey:           have["apikey"],
		IdentityEndpoint: have["authurl"],
	}
	return ao, have["region"], nil
}
