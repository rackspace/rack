package blockstoragecommands

import (
	"github.com/jrperritt/rack/commands/blockstoragecommands/snapshotcommands"
	"github.com/jrperritt/rack/commands/blockstoragecommands/volumecommands"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
)

// Get returns all the commands allowed for a `block-storage` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "snapshot",
			Usage:       "Copies of block storage volumes at a specific moment in time. Use for backups, restores, and long term storage.",
			Subcommands: snapshotcommands.Get(),
		},
		{
			Name:        "volume",
			Usage:       "Block level volumes to expand storage on your servers.",
			Subcommands: volumecommands.Get(),
		},
	}
}
