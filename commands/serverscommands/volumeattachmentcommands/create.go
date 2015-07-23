package volumeattachmentcommands

import (
	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/fatih/structs"
	osVolumeAttach "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
	osServers "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/jrperritt/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "[--server-id <serverID> | --server-name <serverName>] [--id <volumeID> | --name <volumeName> | --stdin id]"),
	Description: "Creates a new volume attachment on the server",
	Action:      actionCreate,
	Flags:       commandoptions.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the volume to attach.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
		cli.StringFlag{
			Name:  "server-id",
			Usage: "[optional; required if `server-name` isn't provided] The server ID to which attach the volume.",
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "[optional; required if `server-id` isn't provided] The server name to which attach the volume.",
		},
		cli.StringFlag{
			Name:  "device",
			Usage: "[optional] The name of the device to which the volume will attach. Default is 'auto'.",
		},
	}
}

var keysCreate = []string{"ID", "Device", "VolumeID", "ServerID"}

type paramsCreate struct {
	opts     *osVolumeAttach.CreateOpts
	serverID string
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
	serverID, err := command.Ctx.IDOrName(osServers.IDFromName)
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext
	opts := &osVolumeAttach.CreateOpts{
		Device: c.String("device"),
	}

	resource.Params = &paramsCreate{
		opts:     opts,
		serverID: serverID,
	}
	return nil
}

func (command *commandCreate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsCreate).opts.VolumeID = item
	return nil
}

func (command *commandCreate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}

	resource.Params.(*paramsCreate).opts.VolumeID = command.Ctx.CLIContext.String("id")
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsCreate)
	volumeAttachment, err := osVolumeAttach.Create(command.Ctx.ServiceClient, params.serverID, params.opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(volumeAttachment)
}

func (command *commandCreate) StdinField() string {
	return "id"
}
