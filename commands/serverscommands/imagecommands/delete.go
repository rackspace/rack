package imagecommands

import (
	"fmt"
	"time"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	osImages "github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/rack/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description: "Deletes an image.",
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
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the image to delete",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the image to delete",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped to STDIN. Valid values are: id",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the image has been deleted",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	imageID string
	wait    bool
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
	wait := false
	if command.Ctx.CLIContext.IsSet("wait-for-completion") {
		wait = true
	}

	resource.Params = &paramsDelete{
		wait: wait,
	}
	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsDelete).imageID = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	id, err := command.Ctx.IDOrName(osImages.IDFromName)
	resource.Params.(*paramsDelete).imageID = id
	return err
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDelete)
	imageID := params.imageID
	err := osImages.Delete(command.Ctx.ServiceClient, imageID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}

	if params.wait {
		i := 0
		for i < 120 {
			_, err := images.Get(command.Ctx.ServiceClient, imageID).Extract()
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
			i++
		}
		resource.Result = fmt.Sprintf("Deleted image [%s]", imageID)
	} else {
		resource.Result = fmt.Sprintf("Deleting instance [%s]\n", imageID)
	}
}

func (command *commandDelete) StdinField() string {
	return "id"
}
