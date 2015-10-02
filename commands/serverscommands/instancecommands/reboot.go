package instancecommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osServers "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/rackspace/rack/output"
	"github.com/rackspace/rack/util"
)

var reboot = cli.Command{
	Name:        "reboot",
	Usage:       util.Usage(commandPrefix, "reboot", "[--id <serverID> | --name <serverName> | --stdin id] [--soft | --hard]"),
	Description: "Reboots an existing server",
	Action:      actionReboot,
	Flags:       commandoptions.CommandFlags(flagsReboot, keysReboot),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsReboot, keysReboot))
	},
}

func flagsReboot() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "soft",
			Usage: "[optional; required if 'hard' is not provided] Ask the OS to restart under its own procedures.",
		},
		cli.BoolFlag{
			Name:  "hard",
			Usage: "[optional; required if 'soft' is not provided] Physically cut power to the machine and then restore it after a brief while.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the server.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the server.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the instance has been rebooted.",
		},
	}
}

var keysReboot = []string{}

type paramsReboot struct {
	wait     bool
	serverID string
	how      osServers.RebootMethod
}

type commandReboot handler.Command

func actionReboot(c *cli.Context) {
	command := &commandReboot{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandReboot) Context() *handler.Context {
	return command.Ctx
}

func (command *commandReboot) Keys() []string {
	return keysReboot
}

func (command *commandReboot) ServiceClientType() string {
	return serviceClientType
}

func (command *commandReboot) HandleFlags(resource *handler.Resource) error {
	c := command.Context().CLIContext
	wait := false
	if c.IsSet("wait-for-completion") {
		wait = true
	}

	var how osServers.RebootMethod
	if c.IsSet("soft") {
		how = osServers.OSReboot
	}
	if c.IsSet("hard") {
		how = osServers.PowerCycle
	}
	if how == "" {
		return output.ErrMissingFlag{Msg: "One of either --soft or --hard must be provided."}
	}
	resource.Params = &paramsReboot{
		how:  how,
		wait: wait,
	}
	return nil
}

func (command *commandReboot) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsReboot).serverID = item
	return nil
}

func (command *commandReboot) HandleSingle(resource *handler.Resource) error {
	id, err := command.Context().IDOrName(osServers.IDFromName)
	resource.Params.(*paramsReboot).serverID = id
	return err
}

func (command *commandReboot) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsReboot)
	serverID := params.serverID
	err := servers.Reboot(command.Context().ServiceClient, serverID, params.how).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}

	if params.wait {
		err = osServers.WaitForStatus(command.Ctx.ServiceClient, serverID, "ACTIVE", 600)
		if err != nil {
			resource.Err = err
			return
		}

		resource.Result = fmt.Sprintf("Rebooted instance [%s]\n", serverID)
	} else {
		resource.Result = fmt.Sprintf("Rebooting instance [%s]\n", serverID)
	}
}

func (command *commandReboot) StdinField() string {
	return "id"
}
