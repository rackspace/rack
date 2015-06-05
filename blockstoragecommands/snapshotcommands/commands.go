package snapshotcommands

import "github.com/codegangsta/cli"

var commandPrefix = "block_storage snapshots"

// Get returns all the commands allowed for a `block_storage snapshots` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		create,
		update,
	}
}
