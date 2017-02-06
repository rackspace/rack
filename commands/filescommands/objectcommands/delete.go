package objectcommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/rack/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "--container <containerName> [--name <objectName> | --stdin name]"),
	Description: "Deletes an object",
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
			Usage: "[required] The name of the container",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the object",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped to STDIN. Valid values are: name",
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

	c := command.Ctx.CLIContext
	containerName := c.String("container")

	if err := CheckContainerExists(command.Ctx.ServiceClient, containerName); err != nil {
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
	containerName := resource.Params.(*paramsDelete).container
	objectName := resource.Params.(*paramsDelete).object
	rawResponse := objects.Delete(command.Ctx.ServiceClient, containerName, objectName, nil)
	if rawResponse.Err != nil {
		resource.Err = rawResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted object [%s] from container [%s]\n", objectName, containerName)
}

func (command *commandDelete) StdinField() string {
	return "name"
}
