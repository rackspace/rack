package imagecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osImages "github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/images"
)

var list = cli.Command{
	Name:        "list",
	Usage:       fmt.Sprintf("%s %s list [flags]", util.Name, commandPrefix),
	Description: "Lists images",
	Action:      commandList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Only list images that have this name.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Only list images that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "Start listing flavors at this flavor ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "Only return this many flavors at most.",
		},
	}
}

var keysList = []string{"ID", "Name", "Status", "MinDisk", "MinRAM"}

func commandList(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("compute")
	opts := osImages.ListOpts{
		Name:   c.String("name"),
		Status: c.String("status"),
		Marker: c.String("marker"),
		Limit:  c.Int("limit"),
	}
	allPages, err := images.ListDetail(client, opts).AllPages()
	if err != nil {
		fmt.Printf("Error listing images: %s\n", err)
		os.Exit(1)
	}
	o, err := images.ExtractImages(allPages)
	if err != nil {
		fmt.Printf("Error listing images: %s\n", err)
		os.Exit(1)
	}

	f := func() interface{} {
		m := make([]map[string]interface{}, len(o))
		for j, image := range o {
			m[j] = structs.Map(image)
		}
		return m
	}
	output.Print(c, &f, keysList)
}
