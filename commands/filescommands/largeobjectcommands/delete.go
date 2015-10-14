package largeobjectcommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/commands/filescommands/objectcommands"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osObjects "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/rack/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "--container <containerName> --name <objectName>"),
	Description: "Deletes a large object",
	Action:      actionDelete,
	Flags:       commandoptions.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name the object to delete.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped to STDIN, if any. Valid values are: name.",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	container string
	object    string
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
	err := command.Ctx.CheckFlagsSet([]string{"container"})
	if err != nil {
		return err
	}

	containerName := command.Ctx.CLIContext.String("container")

	if err := objectcommands.CheckContainerExists(command.Ctx.ServiceClient, containerName); err != nil {
		return err
	}

	resource.Params = &paramsDelete{
		container: containerName,
	}

	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsDelete).object = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).object = command.Ctx.CLIContext.String("name")

	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDelete)
	containerName := params.container
	objectName := params.object

	listOpts := osObjects.ListOpts{
		Prefix: objectName,
	}
	allPages, err := osObjects.List(command.Ctx.ServiceClient, containerName, listOpts).AllPages()
	if err != nil {
		resource.Err = err
		return
	}

	objectNames, err := osObjects.ExtractNames(allPages)
	if err != nil {
		resource.Err = err
		return
	}

	for _, thisName := range objectNames {
		rawResponse := osObjects.Delete(command.Ctx.ServiceClient, containerName, thisName, osObjects.DeleteOpts{})
		if rawResponse.Err != nil {
			resource.Err = rawResponse.Err
			return
		}
	}

	resource.Result = fmt.Sprintf("Deleted object [%s] from container [%s]", objectName, containerName)
}

func (command *commandDelete) StdinField() string {
	return "name"
}
