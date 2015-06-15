package auth

import (
	"os"

	"github.com/rackspace/gophercloud"
)

func envvars() (gophercloud.AuthOptions, string, error) {
	ao := gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("RS_AUTH_URL"),
		Username:         os.Getenv("RS_USERNAME"),
		APIKey:           os.Getenv("RS_API_KEY"),
	}
	region := os.Getenv("RS_REGION_NAME")
	return ao, region, nil
}
