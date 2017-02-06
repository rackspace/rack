package securitygrouprulecommands

import "github.com/codegangsta/cli"

var commandPrefix = "networks security-group-rule"
var serviceClientType = "network"

// Get returns all the commands allowed for a `networks securitygrouprule` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		get,
		remove,
		list,
	}
}
