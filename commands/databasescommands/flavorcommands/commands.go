package flavorcommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "databases flavor"
var serviceClientType = "databases"

// Get returns all the commands allowed for a `databases flavor` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
	}
}
