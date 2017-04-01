package usercommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "databases user"
var serviceClientType = "databases"

// Get returns all the commands allowed for a `databases user` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		list,
		get,
		update,
		remove,
		listAccess,
		grantAccess,
		revokeAccess,
	}
}
