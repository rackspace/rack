package stackcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

// delete is a reserved word in Go.
var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <stackID> | --name <stackName>]"),
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
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the server.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the server.",
		},
	}
}

var keysDelete = keysList

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
	// the behavior of the python-heatclient is to show a list of stacks as the
	// output of stack-delete.
	result, err := stackList(command.Ctx.ServiceClient)
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = result
}
