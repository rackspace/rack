package flavorcommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osFlavors "github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/flavors"
)

var list = cli.Command{
	Name:        "list",
	Usage:       fmt.Sprintf("%s %s list [flags]", util.Name, commandPrefix),
	Description: "Lists flavors",
	Action:      commandList,
	Flags:       util.CommandFlags(flagsList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList))
	},
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

	f := func() interface{} {
		m := make([]map[string]interface{}, len(o))
		for j, flavor := range o {
			m[j] = structs.Map(flavor)
		}
		return m
	}
	keys := []string{"ID", "Name", "RAM", "Disk", "Swap", "VCPUs", "RxTxFactor"}
	output.Print(c, &f, keys)
}
