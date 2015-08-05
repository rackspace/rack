package flavorcommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "servers flavor"
var serviceClientType = "compute"

// Get returns all the commands allowed for a `compute flavors` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
	}
}
