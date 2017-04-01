package instancecommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "databases instance"
var serviceClientType = "databases"

// Get returns all the commands allowed for a `servers instance` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		list,
		get,
		remove,
		getConfig,
		hasRoot,
		enableRoot,
		restart,
		resizeFlavor,
		resizeDisk,
	}
}
