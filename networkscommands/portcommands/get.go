package portcommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osPorts "github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/rackspace/networking/v2/ports"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", ""),
	Description: "Gets a port",
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
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the port.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the port.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysGet = []string{"ID", "Name", "NetworkID", "Status", "MACAddress", "DeviceID", "DeviceOwner", "Up", "FixedIPs", "SecurityGroups", "TenantID"}

type paramsGet struct {
	portID string
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
	resource.Params.(*paramsGet).portID = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	portID, err := command.Ctx.IDOrName(osPorts.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).portID = portID
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	portID := resource.Params.(*paramsGet).portID
	port, err := ports.Get(command.Ctx.ServiceClient, portID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = portSingle(port)
}

func (command *commandGet) StdinField() string {
	return "id"
}
