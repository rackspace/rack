package servercommands

import "github.com/codegangsta/cli"

var commandPrefix = "compute servers"

// Get returns all the commands allowed for a `compute servers` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		create,
		get,
		update,
		remove,
		reboot,
	}
}
