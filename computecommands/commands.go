package computecommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rackcli/computecommands/flavorcommands"
	"github.com/jrperritt/rackcli/computecommands/imagecommands"
	"github.com/jrperritt/rackcli/computecommands/servercommands"
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
	}
}
