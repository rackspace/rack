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
			Usage:       "Used for Cloud Networks network operations",
			Subcommands: networkcommands.Get(),
		},
		{
			Name:        "subnet",
			Usage:       "Used for Cloud Networks subnet operations",
			Subcommands: subnetcommands.Get(),
		},
		{
			Name:        "port",
			Usage:       "Used for Cloud Networks port operations",
			Subcommands: portcommands.Get(),
		},
		{
			Name:        "security-group",
			Usage:       "Used for Cloud Networks security group operations",
			Subcommands: securitygroupcommands.Get(),
		},
		{
			Name:        "security-group-rule",
			Usage:       "Used for Cloud Networks security group rule operations",
			Subcommands: securitygrouprulecommands.Get(),
		},
	}
}
