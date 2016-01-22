package databasecommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "databases database"
var serviceClientType = "databases"

// Get returns all the commands allowed for a `servers instance` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		list,
		remove,
	}
}
