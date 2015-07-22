package containercommands

import "github.com/jrperritt/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "files container"
var serviceClientType = "object-store"

// Get returns all the commands allowed for a `files container` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		list,
		get,
		remove,
		update,
		empty,
	}
}
