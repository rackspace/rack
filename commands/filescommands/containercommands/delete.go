package containercommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	"github.com/rackspace/rack/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--name <containerName> | --stdin name]"),
	Description: "Deletes a container",
	Action:      actionDelete,
	Flags:       commandoptions.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container",
		},
		cli.BoolFlag{
			Name:  "purge",
			Usage: "[optional] If set, this command will delete all objects in the container, and then delete the container.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
		cli.IntFlag{
			Name:  "concurrency",
			Usage: "[optional] If `purge` is provided, the amount of concurrent workers that will empty the container.",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "[optional] By default, every object deleted will be outputted. If --quiet is provided, only a final summary will be outputted.",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	container   string
	purge       bool
	quiet       bool
	concurrency int
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
	resource.Params = &paramsDelete{
		purge:       command.Ctx.CLIContext.Bool("purge"),
		quiet:       command.Ctx.CLIContext.Bool("quiet"),
		concurrency: command.Ctx.CLIContext.Int("concurrency"),
	}
	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsDelete).container = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).container = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDelete)
	containerName := params.container
	if params.purge {
		emptyParams := &handleEmptyParams{
			container:   containerName,
			quiet:       params.quiet,
			concurrency: params.concurrency,
		}
		handleEmpty(command, resource, emptyParams)
	}

	rawResponse := containers.Delete(command.Ctx.ServiceClient, containerName)
	if rawResponse.Err != nil {
		resource.Err = rawResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted container [%s]\n", containerName)
}

func (command *commandDelete) StdinField() string {
	return "name"
}
