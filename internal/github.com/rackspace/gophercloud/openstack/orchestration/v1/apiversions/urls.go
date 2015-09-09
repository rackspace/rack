package apiversions

import "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"

func apiVersionsURL(c *gophercloud.ServiceClient) string {
	return c.Endpoint
}
