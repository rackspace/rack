package networkcommands

import (
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osNetworks "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/networking/v2/networks"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", ""),
	Description: "Creates a new networks instance",
	Action:      actionCreate,
	Flags:       util.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] The name that the network should have.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional] The field being piped into STDIN. Valid values are: name",
		},
		cli.BoolFlag{
			Name:  "up",
			Usage: "[optional] If provided, the network will be up.",
		},
		cli.BoolFlag{
			Name:  "shared",
			Usage: "[optional] If provided, the network is shared among all tenants.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] The ID of the tenant who should own this network.",
		},
	}
}

var keysCreate = []string{"ID", "Name", "Up", "Status", "Shared", "Tenant ID"}

type paramsCreate struct {
	opts *osNetworks.CreateOpts
}

type commandCreate handler.Command

func actionCreate(c *cli.Context) {
	command := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandCreate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandCreate) Keys() []string {
	return keysCreate
}

func (command *commandCreate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	opts := &osNetworks.CreateOpts{
		TenantID: c.String("tenant-id"),
	}
	if c.IsSet("up") {
		up := true
		opts.AdminStateUp = &up
	}
	if c.IsSet("shared") {
		shared := true
		opts.Shared = &shared
	}
	resource.Params = &paramsCreate{
		opts: opts,
	}
	return nil
}

func (command *commandCreate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsCreate).opts.Name = item
	return nil
}

func (command *commandCreate) HandleSingle(resource *handler.Resource) error {
	name := command.Ctx.CLIContext.String("name")
	resource.Params.(*paramsCreate).opts.Name = name
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	network, err := networks.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = networkSingle(network)
}

func (command *commandCreate) StdinField() string {
	return "name"
}
