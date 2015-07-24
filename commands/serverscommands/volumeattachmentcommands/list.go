package volumeattachmentcommands

import (
	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/fatih/structs"
	osVolumeAttach "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/pagination"
	"github.com/jrperritt/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", "[--server-id <serverID> | --server-name <serverName> | --stdin id]"),
	Description: "Lists attachments for the given server",
	Action:      actionList,
	Flags:       commandoptions.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "server-id",
			Usage: "[optional; required if `server-name` or `stdin` isn't provided] The server ID of the attachment.",
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "[optional; required if `server-id` or `stdin` isn't provided] The server name of the attachment.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `server-id` or `server-name` isn't provided] The field being piped into STDIN. Valid values are: server-id",
		},
	}
}

var keysList = []string{"ID", "Device", "VolumeID", "ServerID"}

type paramsList struct {
	serverID string
}

type commandList handler.Command

func actionList(c *cli.Context) {
	command := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandList) Context() *handler.Context {
	return command.Ctx
}

func (command *commandList) Keys() []string {
	return keysList
}

func (command *commandList) ServiceClientType() string {
	return serviceClientType
}

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsList{}
	return nil
}

func (command *commandList) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsList).serverID = item
	return nil
}

func (command *commandList) HandleSingle(resource *handler.Resource) error {
	serverID, err := serverIDorName(command.Ctx)
	if err != nil {
		return err
	}

	resource.Params.(*paramsList).serverID = serverID
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsList)
	err := osVolumeAttach.List(command.Ctx.ServiceClient, params.serverID).EachPage(func(page pagination.Page) (bool, error) {
		volumeAttachments, err := osVolumeAttach.ExtractVolumeAttachments(page)
		if err != nil {
			return false, err
		}
		result := make([]map[string]interface{}, len(volumeAttachments))
		for j, volumeAttachment := range volumeAttachments {
			result[j] = structs.Map(volumeAttachment)
		}
		resource.Result = result
		return false, nil
	})
	if err != nil {
		resource.Err = err
		return
	}
}

func (command *commandList) StdinField() string {
	return "server-id"
}
