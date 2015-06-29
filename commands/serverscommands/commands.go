package serverscommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/commands/serverscommands/flavorcommands"
	"github.com/jrperritt/rack/commands/serverscommands/imagecommands"
	"github.com/jrperritt/rack/commands/serverscommands/instancecommands"
	"github.com/jrperritt/rack/commands/serverscommands/keypaircommands"
)

// Get returns all the commands allowed for a `servers` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "instance",
			Usage:       "Used for Server Instance operations",
			Subcommands: instancecommands.Get(),
		},
		{
			Name:        "image",
			Usage:       "Used for Server Image operations",
			Subcommands: imagecommands.Get(),
		},
		{
			Name:        "flavor",
			Usage:       "Used for Server Flavor operations",
			Subcommands: flavorcommands.Get(),
		},
		{
			Name:        "keypair",
			Usage:       "Used for Server Keypair operations",
			Subcommands: keypaircommands.Get(),
		},
	}
}
