package stackeventcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osStackEvents "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stackevents"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackevents"
	"github.com/rackspace/rack/util"
)

var listResource = cli.Command{
	Name:        "list-resource",
	Usage:       util.Usage(commandPrefix, "list-resource", "[--stack-name <stackName> | --stack-id <stackID>] --resource <resourceName>"),
	Description: "Lists events for a specified stack resource",
	Action:      actionListResource,
	Flags:       commandoptions.CommandFlags(flagsListResource, keysListResource),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsListResource, keysListResource))
	},
}

func flagsListResource() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "stack-name",
			Usage: "[optional; required if `id` isn't specified] The stack name.",
		},
		cli.StringFlag{
			Name:  "stack-id",
			Usage: "[optional; required if `name` isn't specified] The stack id.",
		},
		cli.StringFlag{
			Name:  "resource",
			Usage: "[required] The resource name.",
		},
	}
}

type paramsListResource struct {
	stackName    string
	stackID      string
	resourceName string
}

var keysListResource = []string{"ResourceName", "Time", "ResourceStatusReason", "ResourceStatus", "PhysicalResourceID", "ID"}

type commandListResource handler.Command

func actionListResource(c *cli.Context) {
	command := &commandListResource{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandListResource) Context() *handler.Context {
	return command.Ctx
}

func (command *commandListResource) Keys() []string {
	return keysListResource
}

func (command *commandListResource) ServiceClientType() string {
	return serviceClientType
}

func (command *commandListResource) HandleFlags(resource *handler.Resource) error {
	if err := command.Ctx.CheckFlagsSet([]string{"resource"}); err != nil {
		return err
	}

	c := command.Ctx.CLIContext
	name := c.String("stack-name")
	id := c.String("stack-id")
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params = &paramsListResource{
		stackName:    name,
		stackID:      id,
		resourceName: c.String("resource"),
	}
	return nil
}

func (command *commandListResource) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsListResource)
	stackName := params.stackName
	stackID := params.stackID
	resourceName := params.resourceName

	pager := stackevents.ListResourceEvents(command.Ctx.ServiceClient, stackName, stackID, resourceName, nil)
	pages, err := pager.AllPages()
	if err != nil {
		resource.Err = err
		return
	}
	info, err := osStackEvents.ExtractResourceEvents(pages)
	if err != nil {
		resource.Err = err
		return
	}
	result := make([]map[string]interface{}, len(info))
	for j, event := range info {
		result[j] = structs.Map(&event)
		result[j]["Time"] = event.Time
	}
	resource.Result = result
}
