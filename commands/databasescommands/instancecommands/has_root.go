package instancecommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	"github.com/rackspace/rack/util"
)

var hasRoot = cli.Command{
	Name:        "has-root",
	Usage:       util.Usage(commandPrefix, "has-root", "[--id <serverID> | --stdin id]"),
	Description: "Indicates whether the root user has been enabled on a database instance",
	Action:      actionHasRoot,
	Flags:       commandoptions.CommandFlags(flagsHasRoot, keysHasRoot),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsHasRoot, keysHasRoot))
	},
}

func flagsHasRoot() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the database instance.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
	}
}

var keysHasRoot = []string{}

type paramsHasRoot struct {
	id string
}

type commandHasRoot handler.Command

func actionHasRoot(c *cli.Context) {
	command := &commandHasRoot{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandHasRoot) Context() *handler.Context {
	return command.Ctx
}

func (command *commandHasRoot) Keys() []string {
	return keysHasRoot
}

func (command *commandHasRoot) ServiceClientType() string {
	return serviceClientType
}

func (command *commandHasRoot) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsHasRoot{}
	return nil
}

func (command *commandHasRoot) StdinField() string {
	return "id"
}

func (command *commandHasRoot) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsHasRoot).id = item
	return nil
}

func (command *commandHasRoot) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsHasRoot).id = command.Ctx.CLIContext.String("id")
	return err
}

func (command *commandHasRoot) Execute(resource *handler.Resource) {
	id := resource.Params.(*paramsHasRoot).id
	isEnabled, err := instances.IsRootEnabled(command.Context().ServiceClient, id)
	if err != nil {
		resource.Err = err
		return
	}

	msg := fmt.Sprintln("Root user is enabled")
	if !isEnabled {
		msg = fmt.Sprintln("Root user is not enabled")
	}

	resource.Result = msg
}
