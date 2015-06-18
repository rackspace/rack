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

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", util.IDOrNameUsage("image")),
	Description: "Retreives an image",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return util.IDAndNameFlags
}

var keysGet = []string{"ID", "Name", "Status", "Progress", "MinDisk", "MinRAM", "Created", "Updated"}

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

	imageID, err := util.IDOrName(c, client, osImages.IDFromName)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}
	o, err := images.Get(client, imageID).Extract()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error retrieving image [%s]: %s\n", imageID, err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		return structs.Map(o)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
