package imagecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
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
	Flags:       util.CommandFlags(flagsList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList))
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
	output.Print(c, o, tableList)
}

func tableList(c *cli.Context, i interface{}) {
	images, ok := i.([]osImages.Image)
	if !ok {
		fmt.Fprintf(c.App.Writer, "Could not type assert interface\n%+v\nto []osImages.Image\n", i)
		os.Exit(1)
	}

	keys := []string{"ID", "Name", "Status", "MinDisk", "MinRAM"}

	many := make([]interface{}, len(images))
	for i, d := range images {
		many[i] = d
	}

	util.SimpleListing(c, many, keys)

}
