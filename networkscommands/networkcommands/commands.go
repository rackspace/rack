package networkcommands

import "github.com/codegangsta/cli"

var commandPrefix = "networks network"
var serviceClientType = "network"

// Get returns all the commands allowed for a `files container` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		get,
		remove,
		list,
	}
}
