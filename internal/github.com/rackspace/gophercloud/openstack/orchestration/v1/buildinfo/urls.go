package buildinfo

import "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"

func getURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("build_info")
}
