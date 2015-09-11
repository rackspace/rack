package stackresourcecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackresources"
	"github.com/rackspace/rack/util"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "--name <stackName> --id <stackID> --resource <resourceName>"),
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
			Name:  "name",
			Usage: "[optional; required if `id` isn't provided] The stack name.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `id` isn't provided] The stack id.",
		},
		cli.StringFlag{
			Name:  "resource",
			Usage: "[required] The resource name.",
		},
	}
}

type paramsGet struct {
	stackName    string
	stackID      string
	resourceName string
}

var keysGet = []string{"Name", "PhysicalID", "Type", "Status", "UpdatedTime", "Links", "Attributes", "CreationTime", "Description", "LogicalID"}

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
	if err := command.Ctx.CheckFlagsSet([]string{"resource"}); err != nil {
		return err
	}
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params = &paramsGet{
		stackName:    name,
		stackID:      id,
		resourceName: command.Ctx.CLIContext.String("resource"),
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
	resource.Result = resourceSingle(stackresource)
}
