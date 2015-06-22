package containercommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
)

var get = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "[--container <containerName> | --stdin container]"),
	Description: "Creates a container",
	Action:      actionCreate,
	Flags:       util.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `container` isn't provided] The field being piped into STDIN. Valid values are: container",
		},
	}
}

var keysCreate = []string{}

type paramsCreate struct {
	container string
	opts      containers.CreateOpts
}

type commandCreate handler.Command

func actionCreate(c *cli.Context) {
	command := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandCreate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandCreate) Keys() []string {
	return keysCreate
}

func (command *commandCreate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsCreate{}
	return nil
}

func (command *commandCreate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsCreate).container = item
	return nil
}

func (command *commandCreate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"container"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsCreate).container = command.Ctx.CLIContext.String("container")
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsCreate)
	containerName := params.container
	opts := params.opts
	rawResponse := containers.Create(command.Ctx.ServiceClient, containerName, opts)
	if rawResponse.Err != nil {
		resource.Err = rawResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully created container [%s]\n", containerName)
}

func (command *commandCreate) StdinField() string {
	return "container"
}
