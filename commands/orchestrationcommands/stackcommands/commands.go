package stackcommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "orchestration stack"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for an `orchestration stack` request.
func Get() []cli.Command {
	return []cli.Command{
		abandon,
		adopt,
		create,
		get,
		list,
		listEvents,
		preview,
		remove,
		update,
	}
}
