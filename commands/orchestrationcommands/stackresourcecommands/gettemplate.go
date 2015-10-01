package stackresourcecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackresources"
	"github.com/rackspace/rack/util"
)

var getTemplate = cli.Command{
	Name:        "get-template",
	Usage:       util.Usage(commandPrefix, "get-template", " [--type <resourceType> | --stdin type]"),
	Description: "Show template representation for specified resource type",
	Action:      actionGetTemplate,
	Flags:       commandoptions.CommandFlags(flagsGetTemplate, keysGetTemplate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGetTemplate, keysGetTemplate))
	},
}

func flagsGetTemplate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Usage: "[optional; required if `stdin` isn't provided] The resource type.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `type` isn't provided] The field being piped into STDIN. Valid values are: type.",
		},
	}
}

type paramsGetTemplate struct {
	resourceType string
}

var keysGetTemplate = []string{}

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
	resource.Params.(*paramsGetTemplate).resourceType = item
	return nil
}

func (command *commandGetTemplate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"type"})
	if err != nil {
		return err
	}

	resource.Params = &paramsGetTemplate{
		resourceType: command.Ctx.CLIContext.String("type"),
	}
	return nil
}

func (command *commandGetTemplate) HandleFlags(resource *handler.Resource) error {
	return nil
}

func (command *commandGetTemplate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGetTemplate)
	resourceType := params.resourceType
	template, err := stackresources.Template(command.Ctx.ServiceClient, resourceType).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = string(template)
}

func (command *commandGetTemplate) StdinField() string {
	return "type"
}
