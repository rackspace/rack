package portcommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osPorts "github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/rackspace/networking/v2/ports"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--network-id <network-id>"),
	Description: "Creates a ports",
	Action:      actionCreate,
	Flags:       util.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "network-id",
			Usage: "[required] The network ID of the port.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] The name of the port.",
		},
		cli.BoolFlag{
			Name:  "up",
			Usage: "[optional] If provided, the port will be up upon creation.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] The ID of the tenant that will own the port.",
		},
		cli.StringFlag{
			Name:  "device-id",
			Usage: "[optional] The device ID to associate with the port.",
		},
	}
}

var keysCreate = []string{"ID", "Name", "NetworkID", "Status", "MACAddress", "DeviceID", "DeviceOwner", "Up", "FixedIPs", "SecurityGroups", "TenantID"}

type paramsCreate struct {
	opts *osPorts.CreateOpts
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

	opts := &osPorts.CreateOpts{
		Name:      c.String("name"),
		NetworkID: c.String("network-id"),
		DeviceID:  c.String("device-id"),
		TenantID:  c.String("tenant-id"),
	}

	if c.IsSet("up") {
		up := true
		opts.AdminStateUp = &up
	}

	resource.Params = &paramsCreate{
		opts: opts,
	}

	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	port, err := ports.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = portSingle(port)
}
