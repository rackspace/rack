package volumecommands

import (
	"fmt"

	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osVolumes "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <volumeID> | | --name <volumeName> | --stdin id]"),
	Description: "Deletes a volume",
	Action:      actionDelete,
	Flags:       util.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the volume.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the volume.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	volumeID string
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
	resource.Params = &paramsDelete{}
	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsDelete).volumeID = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	volumeID, err := command.Ctx.IDOrName(osVolumes.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).volumeID = volumeID
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	volumeID := resource.Params.(*paramsDelete).volumeID
	err := osVolumes.Delete(command.Ctx.ServiceClient, volumeID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Deleting volume [%s]\n", volumeID)
}

func (command *commandDelete) StdinField() string {
	return "id"
}
