package snapshotcommands

import "github.com/codegangsta/cli"

var commandPrefix = "blockstorage snapshots"
var serviceClientType = "blockstorage"

// Get returns all the commands allowed for a `blockstorage snapshots` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		create,
	}
}
