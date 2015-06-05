package volumecommands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/olekukonko/tablewriter"
	osVolumes "github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/rackspace/gophercloud/rackspace/blockstorage/v1/volumes"
	"os"
)

var list = cli.Command{
	Name:        "list",
	Usage:       fmt.Sprintf("%s %s list [flags]", util.Name, commandPrefix),
	Description: "Lists volumes",
	Action:      commandList,
}

func commandList(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("blockstorage")
	allPages, err := volumes.List(client).AllPages()
	if err != nil {
		fmt.Printf("Error listing volumes: %s\n", err)
		os.Exit(1)
	}
	o, err := osVolumes.ExtractVolumes(allPages)
	if err != nil {
		fmt.Printf("Error listing volumes: %s\n", err)
		os.Exit(1)
	}
	output.Print(c, o, tableList)
}

func tableList(c *cli.Context, i interface{}) {
	volumes, ok := i.([]osVolumes.Volume)
	if !ok {
		fmt.Fprintf(c.App.Writer, "Could not type assert interface\n%+v\nto []osVolumes.Volume\n", i)
		os.Exit(1)
	}
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	keys := []string{"ID", "Name", "Description", "Size", "VolumeType", "SnapshotID", "Attachments", "CreatedAt"}
	t.SetHeader(keys)
	for _, volume := range volumes {
		m := structs.Map(volume)
		f := []string{}
		for _, key := range keys {
			f = append(f, fmt.Sprintln(m[key]))
		}
		t.Append(f)
	}
	t.Render()
}
