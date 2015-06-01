package auth

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/rackspace"
)

// authentication method via environment variables goes here
func envvars() (gophercloud.AuthOptions, error) {
	return rackspace.AuthOptionsFromEnv()
}
