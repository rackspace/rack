package stackresourcecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	osStackResources "github.com/rackspace/gophercloud/openstack/orchestration/v1/stackresources"
	"github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackresources"
	"github.com/rackspace/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", "[--stack-name <stackName> | --stack-id <stackID> | --stdin stack-name]"),
	Description: "List resources in a stack",
	Action:      actionList,
	Flags:       commandoptions.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
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
		cli.IntFlag{
			Name:  "depth",
			Usage: "[optional] The depth of nested stacks from which resources are retrieved.",
		},
	}
}

type paramsList struct {
	opts      *osStackResources.ListOpts
	stackName string
	stackID   string
}

var keysList = []string{"Name", "PhysicalID", "Type", "Status", "UpdatedTime"}

type commandList handler.Command

func actionList(c *cli.Context) {
	command := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandList) Context() *handler.Context {
	return command.Ctx
}

func (command *commandList) Keys() []string {
	return keysList
}

func (command *commandList) ServiceClientType() string {
	return serviceClientType
}

func (command *commandList) HandlePipe(resource *handler.Resource, item string) error {
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, item, "")
	if err != nil {
		return err
	}
	resource.Params.(*paramsList).stackName = name
	resource.Params.(*paramsList).stackID = id
	return nil
}

func (command *commandList) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("stack-name")
	id := c.String("stack-id")
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params.(*paramsList).stackName = name
	resource.Params.(*paramsList).stackID = id
	return nil
}

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &osStackResources.ListOpts{
		Depth: c.Int("depth"),
	}
	resource.Params = &paramsList{
		opts: opts,
	}
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsList)
	opts := params.opts
	stackName := params.stackName
	stackID := params.stackID
	pager := stackresources.List(command.Ctx.ServiceClient, stackName, stackID, opts)
	pages, err := pager.AllPages()
	if err != nil {
		resource.Err = err
		return
	}
	info, err := osStackResources.ExtractResources(pages)
	if err != nil {
		resource.Err = err
		return
	}
	result := make([]map[string]interface{}, len(info))
	for j, resource := range info {
		result[j] = structs.Map(&resource)
		result[j]["UpdatedTime"] = resource.UpdatedTime
	}
	resource.Result = result
}

func (command *commandList) StdinField() string {
	return "stack-name"
}
