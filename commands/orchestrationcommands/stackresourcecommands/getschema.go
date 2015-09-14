package stackresourcecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackresources"
	"github.com/rackspace/rack/util"
)

var getSchema = cli.Command{
	Name:        "get-schema",
	Usage:       util.Usage(commandPrefix, "get-schema", " [--type <resourceType> | --stdin type]"),
	Description: "Shows the interface schema for a specified resource type.",
	Action:      actionGetSchema,
	Flags:       commandoptions.CommandFlags(flagsGetSchema, keysGetSchema),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGetSchema, keysGetSchema))
	},
}

func flagsGetSchema() []cli.Flag {
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

type paramsGetSchema struct {
	resourceType string
}

var keysGetSchema = []string{"Attributes", "Properties", "ResourceType", "SupportStatus"}

type commandGetSchema handler.Command

func actionGetSchema(c *cli.Context) {
	command := &commandGetSchema{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGetSchema) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGetSchema) Keys() []string {
	return keysGetSchema
}

func (command *commandGetSchema) ServiceClientType() string {
	return serviceClientType
}
func (command *commandGetSchema) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGetSchema).resourceType = item
	return nil
}

func (command *commandGetSchema) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"type"})
	if err != nil {
		return err
	}

	resource.Params = &paramsGetSchema{
		resourceType: command.Ctx.CLIContext.String("type"),
	}
	return nil
}

func (command *commandGetSchema) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGetSchema{}
	return nil
}

func (command *commandGetSchema) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGetSchema)
	resourceType := params.resourceType
	schema, err := stackresources.Schema(command.Ctx.ServiceClient, resourceType).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = resourceSingle(schema)
}

func (command *commandGetSchema) StdinField() string {
	return "type"
}
