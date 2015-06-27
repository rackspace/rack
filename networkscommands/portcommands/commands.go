package portcommands

import "github.com/codegangsta/cli"

var commandPrefix = "networks port"
var serviceClientType = "network"

// Get returns all the commands allowed for a `networks port` request.
func Get() []cli.Command {
	return []cli.Command{
		//create,
		//get,
		//remove,
		list,
		//update,
	}
}
