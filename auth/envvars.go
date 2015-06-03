package auth

import (
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/rackspace"
)

// authentication method via environment variables goes here
func envvars() (gophercloud.AuthOptions, string, error) {
	ao, err := rackspace.AuthOptionsFromEnv()
	region := os.Getenv("RS_REGION_NAME")
	return ao, region, err
}
