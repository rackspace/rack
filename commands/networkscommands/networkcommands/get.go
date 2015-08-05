package networkcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osNetworks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/networking/v2/networks"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", ""),
	Description: "Gets an existing network",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the network",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the network.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysGet = []string{"ID", "Name", "Up", "Status", "Subnets", "Shared", "Tenant ID"}

type paramsGet struct {
	networkID string
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
	resource.Params.(*paramsGet).networkID = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	networkID, err := command.Ctx.IDOrName(osNetworks.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).networkID = networkID
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	networkID := resource.Params.(*paramsGet).networkID
	network, err := networks.Get(command.Ctx.ServiceClient, networkID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = networkSingle(network)
}

func (command *commandGet) StdinField() string {
	return "id"
}
