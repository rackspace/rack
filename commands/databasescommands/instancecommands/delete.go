package instancecommands

import (
	"fmt"
	"time"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	"github.com/rackspace/rack/util"
)

// delete is a reserved word in Go.
var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <instanceId> | --stdin id]"),
	Description: "Deletes an existing database instance",
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
			Usage: "[optional; required if `stdin` isn't provided] The ID of the instance.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the instance has been deleted",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	wait bool
	id   string
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
	resource.Params.(*paramsDelete).id = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).id = command.Ctx.CLIContext.String("id")
	return err
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	id := resource.Params.(*paramsDelete).id
	err := instances.Delete(command.Ctx.ServiceClient, id).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}

	if resource.Params.(*paramsDelete).wait {
		i := 0
		for i < 120 {
			if _, err := instances.Get(command.Ctx.ServiceClient, id).Extract(); err != nil {
				break
			}
			time.Sleep(5 * time.Second)
			i++
		}
		resource.Result = fmt.Sprintf("Deleted instance [%s]\n", id)
	} else {
		resource.Result = fmt.Sprintf("Deleting instance [%s]\n", id)
	}
}

func (command *commandDelete) StdinField() string {
	return "id"
}
