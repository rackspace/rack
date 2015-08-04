package filescommands

import (
	"github.com/jrperritt/rack/commands/filescommands/accountcommands"
	"github.com/jrperritt/rack/commands/filescommands/containercommands"
	"github.com/jrperritt/rack/commands/filescommands/objectcommands"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
)

// Get returns all the commands allowed for a `files` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "account",
			Usage:       "Storage for you account metadata.",
			Subcommands: accountcommands.Get(),
		},
		{
			Name:        "container",
			Usage:       "Storage compartments for your objects/files.",
			Subcommands: containercommands.Get(),
		},
		{
			Name:        "object",
			Usage:       "Data storage for objects/files/media.",
			Subcommands: objectcommands.Get(),
		},
	}
}
