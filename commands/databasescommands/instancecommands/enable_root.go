package instancecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	"github.com/rackspace/rack/util"
)

var enableRoot = cli.Command{
	Name:        "enable-root",
	Usage:       util.Usage(commandPrefix, "enable-root", "[--id <serverID> | --stdin id]"),
	Description: "Enables root user for an existing database instance",
	Action:      actionEnableRoot,
	Flags:       commandoptions.CommandFlags(flagsEnableRoot, keysEnableRoot),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsEnableRoot, keysEnableRoot))
	},
}

func flagsEnableRoot() []cli.Flag {
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

var keysEnableRoot = []string{"Name", "Password"}

type paramsEnableRoot struct {
	id string
}

type commandEnableRoot handler.Command

func actionEnableRoot(c *cli.Context) {
	command := &commandEnableRoot{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandEnableRoot) Context() *handler.Context {
	return command.Ctx
}

func (command *commandEnableRoot) Keys() []string {
	return keysEnableRoot
}

func (command *commandEnableRoot) ServiceClientType() string {
	return serviceClientType
}

func (command *commandEnableRoot) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsEnableRoot{}
	return nil
}

func (command *commandEnableRoot) StdinField() string {
	return "id"
}

func (command *commandEnableRoot) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsEnableRoot).id = item
	return nil
}

func (command *commandEnableRoot) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsEnableRoot).id = command.Ctx.CLIContext.String("id")
	return err
}

func (command *commandEnableRoot) Execute(resource *handler.Resource) {
	id := resource.Params.(*paramsEnableRoot).id
	user, err := instances.EnableRootUser(command.Context().ServiceClient, id).Extract()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = singleUser(user)
}
