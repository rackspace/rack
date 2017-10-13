package stackresourcecommands

import (
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osStackEvents "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stackevents"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackevents"
	"github.com/rackspace/rack/util"
)

var listEvents = cli.Command{
	Name:        "list-events",
	Usage:       util.Usage(commandPrefix, "list-events", "[--stack-name <stackName> | --stack-id <stackID>] --name <resourceName>"),
	Description: "Lists events for a specified stack resource",
	Action:      actionListEvents,
	Flags:       commandoptions.CommandFlags(flagsListEvents, keysListEvents),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsListEvents, keysListEvents))
	},
}

func flagsListEvents() []cli.Flag {
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
			Name:  "name",
			Usage: "[required] The resource name.",
		},
		cli.StringFlag{
			Name:  "resource-actions",
			Usage: "[optional] A comma separated list of actions used for filtering events. Valid values are: CREATE, DELETE, UPDATE, ROLLBACK, SUSPEND, RESUME, ADOPT",
		},
		cli.StringFlag{
			Name:  "resource-statuses",
			Usage: "[optional] A comma separated list of statuses used for filtering events. Valid values are: IN_PROGRESS, COMPLETE, FAILED",
		},
		cli.StringFlag{
			Name:  "resource-types",
			Usage: "[optional] A comma separated list of resource types used for filtering events. e.g. OS::Nova::Server",
		},
		cli.StringFlag{
			Name:  "sort-key",
			Usage: "[optional] Key used to sort the list of stacks. Valid values are: name, status, created_at, updated_at",
		},
		cli.StringFlag{
			Name:  "sort-dir",
			Usage: "[optional] Specify direction for sort. Valid values are: asc, desc",
		},
	}
}

type paramsListEvents struct {
	opts         *osStackEvents.ListResourceEventsOpts
	stackName    string
	stackID      string
	resourceName string
}

var keysListEvents = []string{"ResourceName", "Time", "ResourceStatusReason", "ResourceStatus", "PhysicalResourceID", "ID"}

type commandListEvents handler.Command

func actionListEvents(c *cli.Context) {
	command := &commandListEvents{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandListEvents) Context() *handler.Context {
	return command.Ctx
}

func (command *commandListEvents) Keys() []string {
	return keysListEvents
}

func (command *commandListEvents) ServiceClientType() string {
	return serviceClientType
}

func (command *commandListEvents) HandleFlags(resource *handler.Resource) error {
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
	opts := &osStackEvents.ListResourceEventsOpts{
		ResourceActions:  strings.Split(c.String("resource-actions"), ","),
		ResourceStatuses: strings.Split(c.String("resource-statuses"), ","),
		ResourceTypes:    strings.Split(c.String("resource-types"), ","),
		SortKey:          osStackEvents.SortKey(c.String("sort-key")),
		SortDir:          osStackEvents.SortDir(c.String("sort-dir")),
	}

	resource.Params = &paramsListEvents{
		opts:         opts,
		stackName:    name,
		stackID:      id,
		resourceName: c.String("name"),
	}
	return nil
}

func (command *commandListEvents) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsListEvents)
	opts := params.opts
	stackName := params.stackName
	stackID := params.stackID
	resourceName := params.resourceName

	pager := stackevents.ListResourceEvents(command.Ctx.ServiceClient, stackName, stackID, resourceName, opts)

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
