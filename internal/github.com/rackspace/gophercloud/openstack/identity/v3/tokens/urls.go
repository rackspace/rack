package tokens

import "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud"

func tokenURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
