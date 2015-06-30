package snapshotcommands

import "github.com/jrperritt/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "block-storage snapshots"
var serviceClientType = "blockstorage"

// Get returns all the commands allowed for a `block-storage snapshots` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		create,
		remove,
	}
}
