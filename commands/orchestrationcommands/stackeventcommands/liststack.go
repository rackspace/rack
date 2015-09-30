package stackeventcommands

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

var listStack = cli.Command{
	Name:        "list-stack",
	Usage:       util.Usage(commandPrefix, "list-stack", "[--stack-name <stackName> | --stack-id <stackID> | --stdin stack-name]"),
	Description: "Lists events for a specified stack",
	Action:      actionListStack,
	Flags:       commandoptions.CommandFlags(flagsListStack, keysListStack),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsListStack, keysListStack))
	},
}

func flagsListStack() []cli.Flag {
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

type paramsListStack struct {
	opts      *osStackEvents.ListOpts
	stackName string
	stackID   string
}

var keysListStack = []string{"ResourceName", "Time", "ResourceStatusReason", "ResourceStatus", "PhysicalResourceID", "ID"}

type commandListStack handler.Command

func actionListStack(c *cli.Context) {
	command := &commandListStack{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandListStack) Context() *handler.Context {
	return command.Ctx
}

func (command *commandListStack) Keys() []string {
	return keysListStack
}

func (command *commandListStack) ServiceClientType() string {
	return serviceClientType
}

func (command *commandListStack) HandleFlags(resource *handler.Resource) error {
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
	resource.Params = &paramsListStack{
		opts: opts,
	}
	return nil
}

func (command *commandListStack) HandlePipe(resource *handler.Resource, item string) error {
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, item, "")
	if err != nil {
		return err
	}
	resource.Params.(*paramsListStack).stackName = name
	resource.Params.(*paramsListStack).stackID = id
	return nil
}

func (command *commandListStack) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("stack-name")
	id := c.String("stack-id")
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params.(*paramsListStack).stackName = name
	resource.Params.(*paramsListStack).stackID = id
	return nil
}

func (command *commandListStack) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsListStack)
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

func (command *commandListStack) StdinField() string {
	return "stack-name"
}
