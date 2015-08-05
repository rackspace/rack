package subnetcommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "networks subnet"
var serviceClientType = "network"

// Get returns all the commands allowed for a `networks subnet` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		get,
		remove,
		list,
		update,
	}
}
