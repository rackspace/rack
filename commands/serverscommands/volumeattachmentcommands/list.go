package volumeattachmentcommands

import (
	"fmt"

	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/fatih/structs"
	osVolumeAttach "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
	osServers "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/pagination"
	"github.com/jrperritt/rack/output"
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
			Usage: "[optional; required if `server-id` or `server-name` isn't provided] The field being piped into STDIN. Valid values are: id",
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
	ctx := command.Ctx
	var serverID string
	if ctx.CLIContext.IsSet("server-id") {
		if ctx.CLIContext.IsSet("server-name") {
			return fmt.Errorf("Only one of either --server-id or --server-name may be provided.")
		}
		serverID = ctx.CLIContext.String("server-id")
	} else if ctx.CLIContext.IsSet("server-name") {
		name := ctx.CLIContext.String("server-name")
		id, err := osServers.IDFromName(ctx.ServiceClient, name)
		if err != nil {
			return fmt.Errorf("Error converting name [%s] to ID: %s", name, err)
		}
		serverID = id
	} else {
		return output.ErrMissingFlag{Msg: "One of either --server-id or --server-name must be provided."}
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
	return "id"
}
