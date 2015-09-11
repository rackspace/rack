package orchestrationcommands

import (
	"github.com/rackspace/rack/commands/orchestrationcommands/buildinfocommands"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackeventcommands"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackresourcecommands"
	"github.com/rackspace/rack/commands/orchestrationcommands/templatecommands"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
)

var serviceClientType = "orchestration"

// Get returns all the commands allowed for a `orchestration` request.
func Get() []cli.Command {
	return []cli.Command{
		buildinfocommands.BuildInfo,
		{
			Name:        "stack",
			Usage:       "Stack management.",
			Subcommands: stackcommands.Get(),
		},
		{
			Name:        "event",
			Usage:       "Stack event queries.",
			Subcommands: stackeventcommands.Get(),
		},
		{
			Name:        "resource",
			Usage:       "Stack resource queries.",
			Subcommands: stackresourcecommands.Get(),
		},
		{
			Name:        "template",
			Usage:       "Stack template queries.",
			Subcommands: templatecommands.Get(),
		},
	}
}
