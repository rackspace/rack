package flavorcommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
	"github.com/olekukonko/tablewriter"
	osFlavors "github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/flavors"
)

var list = cli.Command{
	Name:        "list",
	Usage:       fmt.Sprintf("%s %s list [flags]", util.Name, commandPrefix),
	Description: "Lists flavors",
	Action:      commandList,
	Flags:       flagsList(),
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  "minDisk",
			Usage: "[optional] Only list flavors that have at least this much disk storage (in GB).",
		},
		cli.IntFlag{
			Name:  "minRam",
			Usage: "[optional] Only list flavors that have at least this much RAM (in GB).",
		},

		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing flavors at this flavor ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many flavors at most.",
		},
	}
}

func commandList(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("compute")
	opts := flavors.ListOpts{
		MinDisk: c.Int("minDisk"),
		MinRAM:  c.Int("minRam"),
		Marker:  c.String("marker"),
		Limit:   c.Int("limit"),
	}
	allPages, err := flavors.ListDetail(client, opts).AllPages()
	if err != nil {
		fmt.Printf("Error listing flavors: %s\n", err)
		os.Exit(1)
	}
	o, err := osFlavors.ExtractFlavors(allPages)
	if err != nil {
		fmt.Printf("Error listing flavors: %s\n", err)
		os.Exit(1)
	}
	output.Print(c, o, tableList)
}

func tableList(c *cli.Context, i interface{}) {
	flavors, ok := i.([]osFlavors.Flavor)
	if !ok {
		fmt.Fprintf(c.App.Writer, "Could not type assert interface\n%+v\nto []osFlavors.Flavor\n", i)
		os.Exit(1)
	}
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	keys := []string{"ID", "Name", "RAM", "Disk", "Swap", "VCPUs", "RxTxFactor"}
	t.SetHeader(keys)
	for _, flavor := range flavors {
		m := structs.Map(flavor)
		f := []string{}
		for _, key := range keys {
			f = append(f, fmt.Sprint(m[key]))
		}
		t.Append(f)
	}
	t.Render()
}
