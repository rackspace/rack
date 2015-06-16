package imagecommands

import (
	"fmt"

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
	Usage:       util.Usage(commandPrefix, "list", util.IDOrNameUsage("image")),
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
			Usage: "Start listing images at this image ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "Only return this many images at most.",
		},
	}
}

var keysList = []string{"ID", "Name", "Status", "MinDisk", "MinRAM"}

func commandList(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysList,
	}
	err = util.CheckArgNum(c, 0)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	outputParams.ServiceClientType = serviceClientType
	client, err := auth.NewClient(c, outputParams.ServiceClientType)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	opts := osImages.ListOpts{
		Name:   c.String("name"),
		Status: c.String("status"),
		Marker: c.String("marker"),
		Limit:  c.Int("limit"),
	}
	allPages, err := images.ListDetail(client, opts).AllPages()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error listing images: %s\n", err)
		output.Print(outputParams)
		return
	}
	o, err := osImages.ExtractImages(allPages)
	if err != nil {
		outputParams.Err = fmt.Errorf("Error listing images: %s\n", err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		m := make([]map[string]interface{}, len(o))
		for j, image := range o {
			m[j] = structs.Map(image)
		}
		return m
	}
	outputParams.F = &f
	output.Print(outputParams)
}
