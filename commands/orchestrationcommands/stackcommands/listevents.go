package stackcommands

import (
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	osStackEvents "github.com/rackspace/gophercloud/openstack/orchestration/v1/stackevents"
	"github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackevents"
	"github.com/rackspace/rack/util"
)

var listEvents = cli.Command{
	Name:        "list-events",
	Usage:       util.Usage(commandPrefix, "list-events", "[--name <stackName> | --id <stackID> | --stdin name]"),
	Description: "Lists events for a specified stack",
	Action:      actionListEvents,
	Flags:       commandoptions.CommandFlags(flagsListEvents, keysListEvents),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsListEvents, keysListEvents))
	},
}

func flagsListEvents() []cli.Flag {
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
		cli.StringFlag{
			Name:  "resource-actions",
			Usage: "[optional] A comma seperated list of actions used for filtering events. Valid values are: CREATE, DELETE, UPDATE, ROLLBACK, SUSPEND, RESUME, ADOPT",
		},
		cli.StringFlag{
			Name:  "resource-statuses",
			Usage: "[optional] A comma seperated list of statuses used for filtering events. Valid values are: IN_PROGRESS, COMPLETE, FAILED",
		},
		cli.StringFlag{
			Name:  "resource-names",
			Usage: "[optional] A comma seperated list of resource names used for filtering events.",
		},
		cli.StringFlag{
			Name:  "resource-types",
			Usage: "[optional] A comma seperated list of resource types used for filtering events. e.g. OS::Nova::Server",
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
	opts      *osStackEvents.ListOpts
	stackName string
	stackID   string
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
	c := command.Ctx.CLIContext

	stringResourceActions := strings.Split(c.String("resource-actions"), ",")
	resourceActions := make([]osStackEvents.ResourceAction, len(stringResourceActions))
	for i, resourceAction := range stringResourceActions {
		resourceActions[i] = osStackEvents.ResourceAction(resourceAction)
	}

	stringResourceStatuses := strings.Split(c.String("resource-statuses"), ",")
	resourceStatuses := make([]osStackEvents.ResourceStatus, len(stringResourceStatuses))
	for i, resourceStatus := range stringResourceStatuses {
		resourceStatuses[i] = osStackEvents.ResourceStatus(resourceStatus)
	}

	opts := &osStackEvents.ListOpts{
		ResourceActions:  resourceActions,
		ResourceStatuses: resourceStatuses,
		ResourceNames:    strings.Split(c.String("resource-names"), ","),
		ResourceTypes:    strings.Split(c.String("resource-types"), ","),
		SortKey:          osStackEvents.SortKey(c.String("sort-key")),
		SortDir:          osStackEvents.SortDir(c.String("sort-dir")),
	}
	resource.Params = &paramsListEvents{
		opts: opts,
	}
	return nil
}

func (command *commandListEvents) HandlePipe(resource *handler.Resource, item string) error {
	name, id, err := IDAndName(command.Ctx.ServiceClient, item, "")
	if err != nil {
		return err
	}
	resource.Params.(*paramsListEvents).stackName = name
	resource.Params.(*paramsListEvents).stackID = id
	return nil
}

func (command *commandListEvents) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params.(*paramsListEvents).stackName = name
	resource.Params.(*paramsListEvents).stackID = id
	return nil
}

func (command *commandListEvents) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsListEvents)
	opts := params.opts
	stackName := params.stackName
	stackID := params.stackID

	pager := stackevents.List(command.Ctx.ServiceClient, stackName, stackID, opts)
	pages, err := pager.AllPages()
	if err != nil {
		resource.Err = err
		return
	}
	info, err := osStackEvents.ExtractEvents(pages)
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

func (command *commandListEvents) StdinField() string {
	return "name"
}
