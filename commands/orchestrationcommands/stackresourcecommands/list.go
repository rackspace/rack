package stackresourcecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osStackResources "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stackresources"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackresources"
	"github.com/rackspace/rack/util"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", "--name <stackName> --id <stackID>"),
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
			Name:  "name",
			Usage: "[optional; required if `id` isn't provided] The stack name.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't provided] The stack id.",
		},
	}
}

type paramsList struct {
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

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}

	resource.Params = &paramsList{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsList)
	stackName := params.stackName
	stackID := params.stackID
	pager := stackresources.List(command.Ctx.ServiceClient, stackName, stackID, nil)
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
		// TODO: fix the decoding/parsing to make this work right
		result[j]["UpdatedTime"] = resource.UpdatedTime
	}
	resource.Result = result
}
