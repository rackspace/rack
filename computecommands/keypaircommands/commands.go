package keypaircommands

import "github.com/codegangsta/cli"

var commandPrefix = "servers keypair"

// Get returns all the commands allowed for a `compute keypairs` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		create,
		get,
		remove,
	}
}
