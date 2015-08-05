package imagecommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "servers image"
var serviceClientType = "compute"

// Get returns all the commands allowed for a `compute images` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
	}
}
