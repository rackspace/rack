package usercommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/users"
	"github.com/rackspace/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", "--instance <instanceId>"),
	Description: "Lists existing users on a given instance",
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
			Usage: "[required] The ID of the instance that hosts the users.",
		},
	}
}

var keysList = []string{"Name", "Host", "Databases"}

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
	err := command.Ctx.CheckFlagsSet([]string{"instance"})
	if err != nil {
		return err
	}
	resource.Params = &paramsList{instanceID: command.Ctx.CLIContext.String("instance")}
	return nil
}

func (command *commandList) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsList).instanceID = item
	return nil
}

func (command *commandList) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"instance"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsList).instanceID = command.Ctx.CLIContext.String("instance")
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsList)

	pager := users.List(command.Ctx.ServiceClient, params.instanceID)

	pages, err := pager.AllPages()
	if err != nil {
		resource.Err = err
		return
	}

	info, err := users.ExtractUsers(pages)
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("%#v", info)

	result := make([]map[string]interface{}, len(info))
	for j, user := range info {
		result[j] = singleUser(&user)
	}

	resource.Result = result
}
