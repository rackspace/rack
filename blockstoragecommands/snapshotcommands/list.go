package snapshotcommands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
	"github.com/olekukonko/tablewriter"
	osSnapshots "github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/gophercloud/rackspace/blockstorage/v1/snapshots"
	"os"
)

var list = cli.Command{
	Name:        "list",
	Usage:       fmt.Sprintf("%s %s list [flags]", util.Name, commandPrefix),
	Description: "Lists snapshots",
	Action:      commandList,
}

func commandList(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("blockstorage")
	allPages, err := snapshots.List(client).AllPages()
	if err != nil {
		fmt.Printf("Error listing snapshots: %s\n", err)
		os.Exit(1)
	}
	o, err := osSnapshots.ExtractSnapshots(allPages)
	if err != nil {
		fmt.Printf("Error listing snapshots: %s\n", err)
		os.Exit(1)
	}
	output.Print(c, o, tableList)
}

func tableList(c *cli.Context, i interface{}) {
	snapshots, ok := i.([]osSnapshots.Snapshot)
	if !ok {
		fmt.Fprintf(c.App.Writer, "Could not type assert interface\n%+v\nto []osSnapshots.Snapshot\n", i)
		os.Exit(1)
	}
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	keys := []string{"ID", "Name", "Description", "Metadata", "Size", "Status", "VolumeID"}
	t.SetHeader(keys)
	for _, snapshot := range snapshots {
		m := structs.Map(snapshot)
		f := []string{}
		for _, key := range keys {
			f = append(f, fmt.Sprintln(m[key]))
		}
		t.Append(f)
	}
	t.Render()
}
