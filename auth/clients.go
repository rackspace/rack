package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/rackspace"
)

// NewClient creates and returns a Rackspace client for the given service.
func NewClient(c *cli.Context, t string) *gophercloud.ServiceClient {
	ao, region := authMethod(c)
	region = strings.ToUpper(region)
	ao.AllowReauth = true
	if ao.IdentityEndpoint == "" {
		ao.IdentityEndpoint = rackspace.RackspaceUSIdentity
	}

	var pc *gophercloud.ProviderClient

	cacheKey := fmt.Sprintf("%s,%s,%s,%s", ao.Username, ao.APIKey, ao.IdentityEndpoint, region)
	cachedValue, err := GetCacheValue(cacheKey)
	if err == nil {
		if cachedAo, ok := cachedValue["ao"].(gophercloud.AuthOptions); ok && cachedAo == ao {
			if cachedPc, ok := cachedValue["pc"].(gophercloud.ProviderClient); ok {
				pc = &cachedPc
			}
		}
	}

	if pc == nil {
		pc, err = rackspace.AuthenticatedClient(ao)
		if err != nil {
			fmt.Printf("Error creating ProviderClient: %s\n", err)
			os.Exit(1)
		}
	}

	var sc *gophercloud.ServiceClient
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
		fmt.Printf("Error creating ServiceClient (%s): %s\n", err, t)
		os.Exit(1)
	}
	// sc.UserAgent.Prepend("rack/" + util.Version)
	return sc
}

// authMethod determines the appropriate authentication method for the user.
// It returns a gophercloud.AuthOptions object, the region, and the error.
//
func authMethod(c *cli.Context) (gophercloud.AuthOptions, string) {
	have := make(map[string]string)
	need := map[string]string{
		"username": "",
		"apikey":   "",
		"authurl":  "",
		"region":   "",
	}

	cliopts(c, have, need)
	configfile(c, have, need)
	envvars(have, need)

	ao := gophercloud.AuthOptions{
		Username:         have["username"],
		APIKey:           have["apikey"],
		IdentityEndpoint: have["authurl"],
	}
	return ao, have["region"]
}
