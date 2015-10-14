package largeobjectcommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "files large-object"
var serviceClientType = "object-store"

// Get returns all the commands allowed for a `files large-object` request.
func Get() []cli.Command {
	return []cli.Command{
		upload,
		remove,
		//uploadDir,
	}
}
