package orchestrationcommands

import (
    "github.com/rackspace/rack/commands/orchestrationcommands/buildinfocommands"
    "github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
    "github.com/rackspace/rack/internal/github.com/codegangsta/cli"
)

// Get returns all the commands allowed for a `orchestration` request.
func Get() []cli.Command {
    return []cli.Command{
        {
            Name: "buildinfo",
            Usage: "Build information for heat deployment.",
            Subcommands: buildinfocommands.Get(),
        },
        {
            Name: "stack",
            Usage: "Stack management.",
            Subcommands: stackcommands.Get(),
        },
    }
}
