package securitygroupcommands

import "github.com/jrperritt/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "networks security-group"
var serviceClientType = "network"

// Get returns all the commands allowed for a `networks securitygroup` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		get,
		remove,
		list,
	}
}
