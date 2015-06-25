package networkscommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/networkscommands/networkcommands"
)

// Get returns all the commands allowed for a `networks` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "network",
			Usage:       "Used for Cloud Networks network operations",
			Subcommands: networkcommands.Get(),
		},
	}
}
