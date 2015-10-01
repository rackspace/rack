package stacktemplatecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacktemplates"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--stack-name <stackName> | --stack-id <stackID> | --stdin stack-name]"),
	Description: "Get template for specified stack",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "stack-name",
			Usage: "[optional; required if neither `stack-id` nor `stdin` is provided] The stack name.",
		},
		cli.StringFlag{
			Name:  "stack-id",
			Usage: "[optional; required if neither `stack-name` nor `stdin` is provided] The stack id.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if neither `stack-name` nor `stack-id` is provided] The field being piped into STDIN. Valid values are: stack-name.",
		},
	}
}

type paramsGet struct {
	stackName string
	stackID   string
}

var keysGet = []string{"Description", "HeatTemplateVersion", "Parameters", "Resources", "Outputs"}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, item, "")
	if err != nil {
		return err
	}
	resource.Params = &paramsGet{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("stack-name")
	id := c.String("stack-id")
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}

	resource.Params = &paramsGet{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGet)
	stackName := params.stackName
	stackID := params.stackID
	template, err := stacktemplates.Get(command.Ctx.ServiceClient, stackName, stackID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = string(template)
}

func (command *commandGet) StdinField() string {
	return "stack-name"
}
