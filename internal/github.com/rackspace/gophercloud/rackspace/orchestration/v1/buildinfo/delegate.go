package buildinfo

import (
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
)

// Get retreives build info data for the Heat deployment.
// Get retreives data for the given stack template.
func Get(c *gophercloud.ServiceClient) GetResult {
	var res GetResult
	_, res.Err = c.Get(getURL(c), &res.Body, nil)
	return res
}
