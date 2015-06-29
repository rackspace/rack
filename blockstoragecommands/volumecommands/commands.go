package volumecommands

import "github.com/codegangsta/cli"

var commandPrefix = "blockstorage volume"
var serviceClientType = "blockstorage"

// Get returns all the commands allowed for a `blockstorage volumes` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		create,
		update,
	}
}
