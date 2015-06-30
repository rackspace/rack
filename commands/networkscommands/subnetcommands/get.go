package subnetcommands

import (
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osSubnets "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/networking/v2/subnets"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", ""),
	Description: "Gets an existing subnet",
	Action:      actionGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the subnet",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the subnet.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysGet = []string{"ID", "Name", "Network ID", "CIDR", "EnableDHCP", "Gateway IP", "DNS Nameservers", "Allocation Pools", "Host Routes"}

type paramsGet struct {
	subnetID string
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
	resource.Params.(*paramsGet).subnetID = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	subnetID, err := command.Ctx.IDOrName(osSubnets.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).subnetID = subnetID
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	subnetID := resource.Params.(*paramsGet).subnetID
	subnet, err := subnets.Get(command.Ctx.ServiceClient, subnetID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = subnetSingle(subnet)
}

func (command *commandGet) StdinField() string {
	return "id"
}
