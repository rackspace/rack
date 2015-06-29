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
	Usage:       util.Usage(commandPrefix, "create", "[--name <containerName> | --stdin name]"),
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
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] Comma-separated key-value pairs for the container. Example: key1=val1,key2=val2",
		},
		cli.StringFlag{
			Name:  "container-read",
			Usage: "[optional] Comma-separated list of users for whom to grant read access to the container",
		},
		cli.StringFlag{
			Name:  "container-write",
			Usage: "[optional] Comma-separated list of users for whom to grant write access to the container",
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
	c := command.Ctx.CLIContext
	opts := containers.CreateOpts{
		ContainerRead:  c.String("container-read"),
		ContainerWrite: c.String("container-write"),
	}
	if c.IsSet("metadata") {
		metadata, err := command.Ctx.CheckKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}
	resource.Params = &paramsCreate{
		opts: opts,
	}
	return nil
}

func (command *commandCreate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsCreate).container = item
	return nil
}

func (command *commandCreate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsCreate).container = command.Ctx.CLIContext.String("name")
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
	return "name"
}
