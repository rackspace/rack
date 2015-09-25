package stackcommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "orchestration stack"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for a `orchestartion stack` request.
func Get() []cli.Command {
	return []cli.Command{
		abandon,
		adopt,
		create,
		get,
		list,
		preview,
		remove,
		update,
	}
}
