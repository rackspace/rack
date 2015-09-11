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

var listStack = cli.Command{
	Name:        "list-stack",
	Usage:       util.Usage(commandPrefix, "list-stack", "[--name <stackName> | --id <stackID> | --stdin name]"),
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

type paramsListStack struct {
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
	return nil
}

func (command *commandListStack) HandlePipe(resource *handler.Resource, item string) error {
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, item, "")
	if err != nil {
		return err
	}
	resource.Params = &paramsListStack{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandListStack) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}

	resource.Params = &paramsListStack{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandListStack) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsListStack)
	stackName := params.stackName
	stackID := params.stackID

	pager := stackevents.List(command.Ctx.ServiceClient, stackName, stackID, nil)
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
		// TODO: fix the decoding/parsing to make this work right
		result[j]["Time"] = event.Time
	}
	resource.Result = result
}

func (command *commandListStack) StdinField() string {
	return "name"
}
