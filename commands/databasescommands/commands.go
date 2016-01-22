package databasescommands

import (
	"github.com/rackspace/rack/commands/databasescommands/backupcommands"
	"github.com/rackspace/rack/commands/databasescommands/databasecommands"
	"github.com/rackspace/rack/commands/databasescommands/flavorcommands"
	"github.com/rackspace/rack/commands/databasescommands/instancecommands"
	"github.com/rackspace/rack/commands/databasescommands/usercommands"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
)

// Get returns all the commands allowed for a `databases` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "instance",
			Usage:       "Operations on database instances",
			Subcommands: instancecommands.Get(),
		},
		{
			Name:        "database",
			Usage:       "Operations on databases",
			Subcommands: databasecommands.Get(),
		},
		{
			Name:        "user",
			Usage:       "Operations on users",
			Subcommands: usercommands.Get(),
		},
		{
			Name:        "flavor",
			Usage:       "Operations on flavors",
			Subcommands: flavorcommands.Get(),
		},
		{
			Name:        "backup",
			Usage:       "Operations on backups",
			Subcommands: backupcommands.Get(),
		},
	}
}
