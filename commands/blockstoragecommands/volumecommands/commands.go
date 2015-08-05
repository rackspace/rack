package volumecommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "block-storage volume"
var serviceClientType = "blockstorage"

// Get returns all the commands allowed for a `block-storage volumes` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		create,
		update,
		remove,
	}
}
