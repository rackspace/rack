package volumeattachmentcommands

import (
	"fmt"

	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	osVolumeAttach "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
	osServers "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/jrperritt/rack/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--server-id <serverID> | --server-name <serverName>] --id <attachmentID> "),
	Description: "Deletes an existing volume attachment",
	Action:      actionDelete,
	Flags:       commandoptions.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[required] The ID of the attachment.",
		},
		cli.StringFlag{
			Name:  "server-id",
			Usage: "[optional; required if `server-name` isn't provided] The server ID of the attachment.",
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "[optional; required if `server-id` isn't provided] The server name of the attachment.",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	attachmentID string
	serverID     string
}

type commandDelete handler.Command

func actionDelete(c *cli.Context) {
	command := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandDelete) Context() *handler.Context {
	return command.Ctx
}

func (command *commandDelete) Keys() []string {
	return keysDelete
}

func (command *commandDelete) ServiceClientType() string {
	return serviceClientType
}

func (command *commandDelete) HandleFlags(resource *handler.Resource) error {
	serverID, err := command.Ctx.IDOrName(osServers.IDFromName)
	if err != nil {
		return err
	}

	err = command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}

	resource.Params = &paramsDelete{
		attachmentID: command.Ctx.CLIContext.String("id"),
		serverID:     serverID,
	}
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDelete)
	err := osVolumeAttach.Delete(command.Ctx.ServiceClient, params.serverID, params.attachmentID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted volume attachment [%s]\n", params.attachmentID)
}
