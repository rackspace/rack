package blockstoragecommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/commands/blockstoragecommands/snapshotcommands"
	"github.com/jrperritt/rack/commands/blockstoragecommands/volumecommands"
)

// Get returns all the commands allowed for a `block-storage` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "snapshot",
			Usage:       "Used for BlockStorage Snapshot operations",
			Subcommands: snapshotcommands.Get(),
		},
		{
			Name:        "volume",
			Usage:       "Used for BlockStorage Volume operations",
			Subcommands: volumecommands.Get(),
		},
	}
}
