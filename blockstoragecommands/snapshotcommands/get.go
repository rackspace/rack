package snapshotcommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/blockstorage/v1/snapshots"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get <snapshotID> [flags]", util.Name, commandPrefix),
	Description: "Retreives a snapshot",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{}
}

func commandGet(c *cli.Context) {
	util.CheckArgNum(c, 1)
	snapshotID := c.Args()[0]
	client := auth.NewClient("blockstorage")
	o, err := snapshots.Get(client, snapshotID).Extract()
	if err != nil {
		fmt.Printf("Error retreiving snapshot [%s]: %s\n", snapshotID, err)
		os.Exit(1)
	}
	output.Print(c, o, tableGet)
}

func tableGet(c *cli.Context, i interface{}) {
	keys := []string{"ID", "Name", "Description", "Metadata", "Size", "Status", "VolumeID"}
	output.MetaDataTable(c, i, keys)
}
