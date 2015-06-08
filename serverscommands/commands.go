package computecommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/computecommands/flavorcommands"
	"github.com/jrperritt/rack/computecommands/imagecommands"
	"github.com/jrperritt/rack/computecommands/keypaircommands"
	"github.com/jrperritt/rack/computecommands/servercommands"
)

// Get returns all the commands allowed for a `compute` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "instance",
			Usage:       "Used for Server Instance operations",
			Subcommands: servercommands.Get(),
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
