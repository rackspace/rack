package flavorcommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	osFlavors "github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", util.IDOrNameUsage("flavor")),
	Description: "Retreives a flavor",
	Action:      actionGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	stdin := []cli.Flag{
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped to STDIN. Valid values are: id",
		},
	}
	return append(stdin, util.IDAndNameFlags...)
}

var keysGet = []string{"ID", "Name", "Disk", "RAM", "RxTxFactor", "Swap", "VCPUs"}

type paramsGet struct {
	flavor string
}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGet).flavor = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	id, err := command.Ctx.IDOrName(osFlavors.IDFromName)
	resource.Params.(*paramsGet).flavor = id
	return err
}

func (command *commandGet) Execute(resource *handler.Resource) {
	flavorID := resource.Params.(*paramsGet).flavor
	flavor, err := flavors.Get(command.Ctx.ServiceClient, flavorID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(flavor)
}

func (command *commandGet) StdinField() string {
	return "id"
}
