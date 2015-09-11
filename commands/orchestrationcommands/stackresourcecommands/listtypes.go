package stackresourcecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStackResources "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stackresources"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackresources"
	"github.com/rackspace/rack/util"
	"sort"
)

var listTypes = cli.Command{
	Name:        "list-types",
	Usage:       util.Usage(commandPrefix, "list-types", ""),
	Description: "List all supported template resource types",
	Action:      actionListTypes,
	Flags:       commandoptions.CommandFlags(flagsListTypes, keysListTypes),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsListTypes, keysListTypes))
	},
}

func flagsListTypes() []cli.Flag {
	return []cli.Flag{}
}

var keysListTypes = []string{"ResourceType"}

type commandListTypes handler.Command

func actionListTypes(c *cli.Context) {
	command := &commandListTypes{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandListTypes) Context() *handler.Context {
	return command.Ctx
}

func (command *commandListTypes) Keys() []string {
	return keysListTypes
}

func (command *commandListTypes) ServiceClientType() string {
	return serviceClientType
}

func (command *commandListTypes) HandlePipe(resource *handler.Resource) error {
	return nil
}

func (command *commandListTypes) HandleFlags(resource *handler.Resource) error {
	return nil
}

func (command *commandListTypes) Execute(resource *handler.Resource) {
	pager := stackresources.ListTypes(command.Ctx.ServiceClient)
	pages, err := pager.AllPages()
	if err != nil {
		resource.Err = err
		return
	}
	resource_types, err := osStackResources.ExtractResourceTypes(pages)
	if err != nil {
		resource.Err = err
		return
	}
	sort.Sort(resource_types)
	result := make([]map[string]interface{}, len(resource_types))
	for i, resource_type := range resource_types {
		result[i] = make(map[string]interface{})
		result[i]["ResourceType"] = resource_type
	}
	resource.Result = result
}
