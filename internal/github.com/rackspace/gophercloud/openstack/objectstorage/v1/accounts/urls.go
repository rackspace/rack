package accounts

import "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud"

func getURL(c *gophercloud.ServiceClient) string {
	return c.Endpoint
}

func updateURL(c *gophercloud.ServiceClient) string {
	return getURL(c)
}
