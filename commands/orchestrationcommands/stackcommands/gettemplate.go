package stackcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacktemplates"
	"github.com/rackspace/rack/util"
)

var getTemplate = cli.Command{
	Name:        "get-template",
	Usage:       util.Usage(commandPrefix, "get-template", "[--name <stackName> | --id <stackID> | --stdin name]"),
	Description: "Get template for specified stack",
	Action:      actionGetTemplate,
	Flags:       commandoptions.CommandFlags(flagsGetTemplate, keysGetTemplate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGetTemplate, keysGetTemplate))
	},
}

func flagsGetTemplate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if neither `id` nor `stdin` is provided] The stack name.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if neither `name` nor `stdin` is provided] The stack id.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if neither `name` nor `id` is provided] The field being piped into STDIN. Valid values are: name.",
		},
	}
}

type paramsGetTemplate struct {
	stackName string
	stackID   string
}

var keysGetTemplate = []string{"Description", "HeatTemplateVersion", "Parameters", "Resources", "Outputs"}

type commandGetTemplate handler.Command

func actionGetTemplate(c *cli.Context) {
	command := &commandGetTemplate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGetTemplate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGetTemplate) Keys() []string {
	return keysGetTemplate
}

func (command *commandGetTemplate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGetTemplate) HandlePipe(resource *handler.Resource, item string) error {
	name, id, err := IDAndName(command.Ctx.ServiceClient, item, "")
	if err != nil {
		return err
	}
	resource.Params = &paramsGetTemplate{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandGetTemplate) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}

	resource.Params = &paramsGetTemplate{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandGetTemplate) HandleFlags(resource *handler.Resource) error {
	return nil
}

func (command *commandGetTemplate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGetTemplate)
	stackName := params.stackName
	stackID := params.stackID
	template, err := stacktemplates.Get(command.Ctx.ServiceClient, stackName, stackID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = string(template)
}

func (command *commandGetTemplate) StdinField() string {
	return "name"
}
