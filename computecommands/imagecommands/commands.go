package imagecommands

import "github.com/codegangsta/cli"

var commandPrefix = "compute images"

// Get returns all the commands allowed for a `compute images` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
	}
}
