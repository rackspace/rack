package volumeattachmentcommands

import "github.com/jrperritt/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "servers volume-attachment"
var serviceClientType = "compute"

// Get returns all the commands allowed for a `servers volumeattachment` request.
func Get() []cli.Command {
	return []cli.Command{
		//list,
		create,
		get,
		//remove,
	}
}
