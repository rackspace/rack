package volumecommands

import "github.com/codegangsta/cli"

var commandPrefix = "block_storage volumes"

// Get returns all the commands allowed for a `block_storage volumes` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		create,
		update,
	}
}
