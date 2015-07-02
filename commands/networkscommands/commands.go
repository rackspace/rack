package networkscommands

import (
	"github.com/jrperritt/rack/commands/networkscommands/networkcommands"
	"github.com/jrperritt/rack/commands/networkscommands/portcommands"
	"github.com/jrperritt/rack/commands/networkscommands/securitygroupcommands"
	"github.com/jrperritt/rack/commands/networkscommands/securitygrouprulecommands"
	"github.com/jrperritt/rack/commands/networkscommands/subnetcommands"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
)

// Get returns all the commands allowed for a `networks` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "network",
			Usage:       "Set up new networks to use by servers.",
			Subcommands: networkcommands.Get(),
		},
		{
			Name:        "subnet",
			Usage:       "Allocate IP address blocks, gateways, DNS servers, and host routes to networks.",
			Subcommands: subnetcommands.Get(),
		},
		{
			Name:        "port",
			Usage:       "Virtual switch ports on logical network switches.",
			Subcommands: portcommands.Get(),
		},
		{
			Name:        "security-group",
			Usage:       "Collections of rules for network traffic.",
			Subcommands: securitygroupcommands.Get(),
		},
		{
			Name:        "security-group-rule",
			Usage:       "Rules for network traffic.",
			Subcommands: securitygrouprulecommands.Get(),
		},
	}
}
