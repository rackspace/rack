package networkcommands

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osNetworks "github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/rackspace/networking/v2/networks"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", ""),
	Description: "Updates a new networks instance",
	Action:      actionUpdate,
	Flags:       util.CommandFlags(flagsUpdate, keysUpdate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsUpdate, keysUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the network to update.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the network to update.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] The name that the network should have.",
		},
		cli.StringFlag{
			Name:  "up",
			Usage: "[optional] Whether or not the newtork should be up. Options are: true, false.",
		},
		cli.StringFlag{
			Name:  "shared",
			Usage: "[optional] Whether or not the network should be shared among tenants. Options are: true, false.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] The ID of the tenant who should own this network.",
		},
	}
}

var keysUpdate = []string{"ID", "Name", "Up", "Status", "Shared", "Tenant ID"}

type paramsUpdate struct {
	networkID string
	opts      *osNetworks.UpdateOpts
}

type commandUpdate handler.Command

func actionUpdate(c *cli.Context) {
	command := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUpdate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUpdate) Keys() []string {
	return keysUpdate
}

func (command *commandUpdate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUpdate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	opts := &osNetworks.UpdateOpts{
		TenantID: c.String("tenant-id"),
	}
	if c.IsSet("up") {
		upRaw := c.String("up")
		up, err := strconv.ParseBool(upRaw)
		if err != nil {
			return fmt.Errorf("Invalid value for flag `up`: %s. Options are: true, false", upRaw)
		}
		opts.AdminStateUp = &up
	}
	if c.IsSet("shared") {
		sharedRaw := c.String("shared")
		shared, err := strconv.ParseBool(sharedRaw)
		if err != nil {
			return fmt.Errorf("Invalid value for flag `shared`: %s. Options are: true, false", sharedRaw)
		}
		opts.Shared = &shared
	}
	resource.Params = &paramsUpdate{
		opts: opts,
	}
	return nil
}

func (command *commandUpdate) HandlePipe(resource *handler.Resource, networkID string) error {
	resource.Params.(*paramsUpdate).networkID = networkID
	return nil
}

func (command *commandUpdate) HandleSingle(resource *handler.Resource) error {
	networkID := command.Ctx.CLIContext.String("id")
	resource.Params.(*paramsUpdate).networkID = networkID
	return nil
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	newtorkID := resource.Params.(*paramsUpdate).networkID
	opts := resource.Params.(*paramsUpdate).opts
	network, err := networks.Update(command.Ctx.ServiceClient, newtorkID, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = networkSingle(network)
}

func (command *commandUpdate) StdinField() string {
	return "id"
}
