package containercommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/rackspace/rack/util"
)

var empty = cli.Command{
	Name:        "empty",
	Usage:       util.Usage(commandPrefix, "empty", "[--name <containerName> | --stdin name]"),
	Description: "Deletes all objects in a container, but not the container itself.",
	Action:      actionEmpty,
	Flags:       commandoptions.CommandFlags(flagsEmpty, keysEmpty),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsEmpty, keysEmpty))
	},
}

func flagsEmpty() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
		cli.IntFlag{
			Name:  "concurrency",
			Usage: "[optional] The amount of concurrent workers that will empty the container.",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "[optional] By default, every object deleted will be outputted. If --quiet is provided, only a final summary will be outputted.",
		},
	}
}

var keysEmpty = []string{}

type paramsEmpty struct {
	container   string
	quiet       bool
	concurrency int
}

type commandEmpty handler.Command

func actionEmpty(c *cli.Context) {
	command := &commandEmpty{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandEmpty) Context() *handler.Context {
	return command.Ctx
}

func (command *commandEmpty) Keys() []string {
	return keysEmpty
}

func (command *commandEmpty) ServiceClientType() string {
	return serviceClientType
}

func (command *commandEmpty) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsEmpty{
		quiet:       command.Ctx.CLIContext.Bool("quiet"),
		concurrency: command.Ctx.CLIContext.Int("concurrency"),
	}
	return nil
}

func (command *commandEmpty) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsEmpty).container = item
	return nil
}

func (command *commandEmpty) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsEmpty).container = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandEmpty) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsEmpty)
	emptyParams := &handleEmptyParams{
		concurrency: params.concurrency,
		quiet:       params.quiet,
		container:   params.container,
	}
	handleEmpty(command, resource, emptyParams)
}

func (command *commandEmpty) StdinField() string {
	return "name"
}
