package stackcommands

import (
	"fmt"
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

// delete is a reserved word in Go.
var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--name <stackName> | --id <stackID> | --stdin name]"),
	Description: "Deletes an existing stack",
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
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the stack.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the stack.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: name.",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	stackName string
	stackID   string
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
	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	name, id, err := IDAndName(command.Ctx.ServiceClient, item, "")
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).stackName = name
	resource.Params.(*paramsDelete).stackID = id
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params = &paramsDelete{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDelete)
	stackName := params.stackName
	stackID := params.stackID
	err := stacks.Delete(command.Ctx.ServiceClient, stackName, stackID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Stack %s is being deleted.", stackName)
}

func (command *commandDelete) StdinField() string {
	return "name"
}
