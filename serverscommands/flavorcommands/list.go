package flavorcommands

import (
	"fmt"

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
	Usage:       util.Usage(commandPrefix, "list", util.IDOrNameUsage("flavor")),
	Description: "Lists flavors",
	Action:      commandList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
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

var keysList = []string{"ID", "Name", "RAM", "Disk", "Swap", "VCPUs", "RxTxFactor"}

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

	opts := flavors.ListOpts{
		MinDisk: c.Int("minDisk"),
		MinRAM:  c.Int("minRam"),
		Marker:  c.String("marker"),
		Limit:   c.Int("limit"),
	}
	allPages, err := flavors.ListDetail(client, opts).AllPages()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error listing flavors: %s\n", err)
		output.Print(outputParams)
		return
	}
	o, err := osFlavors.ExtractFlavors(allPages)
	if err != nil {
		outputParams.Err = fmt.Errorf("Error listing flavors: %s\n", err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		m := make([]map[string]interface{}, len(o))
		for j, flavor := range o {
			m[j] = structs.Map(flavor)
		}
		return m
	}
	outputParams.F = &f
	output.Print(outputParams)
}
