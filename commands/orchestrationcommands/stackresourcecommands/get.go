package stackresourcecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackresources"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--stack-name <stackName> | --stack-id <stackID>] --name <resourceName>"),
	Description: "Show data for specified resource",
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
			Usage: "[optional; required if `stack-id` isn't provided] The stack name.",
		},
		cli.StringFlag{
			Name:  "stack-id",
			Usage: "[optional; required if `stack-name` isn't provided] The stack id.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The resource name.",
		},
	}
}

type paramsGet struct {
	stackName    string
	stackID      string
	resourceName string
}

var keysGet = []string{"Attributes", "CreationTime", "Description", "Links", "LogicalID", "Name", "PhysicalID", "RequiredBy", "Status", "StatusReason", "Type", "UpdatedTime"}

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

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	if err := command.Ctx.CheckFlagsSet([]string{"name"}); err != nil {
		return err
	}
	c := command.Ctx.CLIContext
	name := c.String("stack-name")
	id := c.String("stack-id")
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params = &paramsGet{
		stackName:    name,
		stackID:      id,
		resourceName: command.Ctx.CLIContext.String("name"),
	}
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGet)
	stackName := params.stackName
	stackID := params.stackID
	resourceName := params.resourceName
	stackresource, err := stackresources.Get(command.Ctx.ServiceClient, stackName, stackID, resourceName).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = stackresource
}

func (command *commandGet) PreCSV(resource *handler.Resource) error {
	resource.Result = resourceSingle(resource.Result)
	resource.FlattenMap("Attributes")
	resource.FlattenMap("Links")
	resource.FlattenMap("RequiredBy")
	return nil
}

func (command *commandGet) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
