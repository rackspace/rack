package databasecommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing databases on a given instance",
	Action:      actionList,
	Flags:       commandoptions.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "instance",
			Usage: "[required] The ID of the instance that hosts the databases.",
		},
	}
}

var keysList = []string{"Name"}

type paramsList struct {
	instanceID string
}

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
	resource.Params = &paramsList{
		instanceID: command.Ctx.CLIContext.String("instance"),
	}
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsList)
	pager := databases.List(command.Ctx.ServiceClient, params.instanceID)

	pages, err := pager.AllPages()
	if err != nil {
		resource.Err = err
		return
	}

	info, err := databases.ExtractDBs(pages)
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("%#v", info)

	result := make([]map[string]interface{}, len(info))
	for j, database := range info {
		result[j] = map[string]interface{}{"Name": database.Name}
	}

	resource.Result = result
}
