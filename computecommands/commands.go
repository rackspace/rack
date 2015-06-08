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
			Name:        "servers",
			Usage:       "Used for Compute Server operations",
			Subcommands: servercommands.Get(),
		},
		{
			Name:        "images",
			Usage:       "Used for Compute Image operations",
			Subcommands: imagecommands.Get(),
		},
		{
			Name:        "flavors",
			Usage:       "Used for Compute Flavor operations",
			Subcommands: flavorcommands.Get(),
		},
		{
			Name:        "keypairs",
			Usage:       "Used for Compute Keypair operations",
			Subcommands: keypaircommands.Get(),
		},
	}
}
