package stackeventcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackevents"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--stack-name <stackName> | --stack-id <stackID>] --resource <resourceName> --id <eventID>"),
	Description: "Show details for a specified event",
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
			Usage: "[optional; required if `stack-id` isn't specified] The stack name.",
		},
		cli.StringFlag{
			Name:  "stack-id",
			Usage: "[optional; required if `stack-name` isn't specified] The stack id.",
		},
		cli.StringFlag{
			Name:  "resource",
			Usage: "[required] The resource name.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[required] The event id.",
		},
	}
}

type paramsGet struct {
	stackName    string
	stackID      string
	resourceName string
	eventID      string
}

var keysGet = []string{"ResourceName", "Time", "ResourceStatusReason", "Links", "LogicalResourceID", "ResourceStatusReason", "ResourceStatus", "PhysicalResourceID", "ID", "ResourceProperties"}

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
	err := command.Ctx.CheckFlagsSet([]string{"resource", "id"})
	if err != nil {
		return err
	}
	c := command.Ctx.CLIContext
	name := c.String("stack-name")
	id := c.String("stack-id")
	name, id, err = stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params = &paramsGet{
		stackName:    name,
		stackID:      id,
		resourceName: c.String("resource"),
		eventID:      c.String("id"),
	}
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGet)
	stackName := params.stackName
	stackID := params.stackID
	resourceName := params.resourceName
	eventID := params.eventID

	event, err := stackevents.Get(command.Ctx.ServiceClient, stackName, stackID, resourceName, eventID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = event
}

func (command *commandGet) PreCSV(resource *handler.Resource) error {
	resource.Result = eventSingle(resource.Result)
	resource.FlattenMap("Links")
	resource.FlattenMap("ResourceProperties")
	return nil
}

func (command *commandGet) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
