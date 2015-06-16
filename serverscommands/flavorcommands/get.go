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

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", util.IDOrNameUsage("flavor")),
	Description: "Retreives a flavor",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return util.IDAndNameFlags
}

var keysGet = []string{"ID", "Name", "Disk", "RAM", "RxTxFactor", "Swap", "VCPUs"}

func commandGet(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysGet,
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

	flavorID, err := util.IDOrName(c, client, osFlavors.IDFromName)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}
	o, err := flavors.Get(client, flavorID).Extract()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error retrieving flavor [%s]: %s\n", flavorID, err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		return structs.Map(o)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
