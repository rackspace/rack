package portcommands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osPorts "github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/rackspace/networking/v2/ports"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", ""),
	Description: "Updates a ports",
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
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the port to update.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the port to update.",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] A new name for the port.",
		},
		cli.StringFlag{
			Name:  "up",
			Usage: "[optional] Whether or not the port is up. Options are: true, false.",
		},
		cli.StringFlag{
			Name:  "security-groups",
			Usage: "[optional] A comma-separated list of security group IDs for this port.",
		},
		cli.StringFlag{
			Name:  "device-id",
			Usage: "[optional] A device ID to associate with the port.",
		},
	}
}

var keysUpdate = []string{"ID", "Name", "NetworkID", "Status", "MACAddress", "DeviceID", "DeviceOwner", "Up", "FixedIPs", "SecurityGroups", "TenantID"}

type paramsUpdate struct {
	portID string
	opts   *osPorts.UpdateOpts
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
	portID, err := command.Ctx.IDOrName(osPorts.IDFromName)
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext
	opts := &osPorts.UpdateOpts{
		Name:     c.String("rename"),
		DeviceID: c.String("device-id"),
	}

	if c.IsSet("up") {
		upRaw := c.String("up")
		up, err := strconv.ParseBool(upRaw)
		if err != nil {
			return fmt.Errorf("Invalid value for flag `up`: %s. Options are: true, false", upRaw)
		}
		opts.AdminStateUp = &up
	}

	if c.IsSet("security-groups") {
		opts.SecurityGroups = strings.Split(c.String("security-groups"), ",")
	}

	resource.Params = &paramsUpdate{
		portID: portID,
		opts:   opts,
	}

	return nil
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsUpdate).opts
	portID := resource.Params.(*paramsUpdate).portID
	port, err := ports.Update(command.Ctx.ServiceClient, portID, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = portSingle(port)
}
