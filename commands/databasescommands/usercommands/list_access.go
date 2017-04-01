package usercommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/users"
	"github.com/rackspace/rack/util"
)

var listAccess = cli.Command{
	Name:        "list-access",
	Usage:       util.Usage(commandPrefix, "list-access", "--instance <instanceId> --name <username>"),
	Description: "Shows a list of all databases a user has access to.",
	Action:      actionListAccess,
	Flags:       commandoptions.CommandFlags(flagsListAccess, keysListAccess),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsListAccess, keysListAccess))
	},
}

func flagsListAccess() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "instance",
			Usage: "[required] The ID of the instance that hosts the users.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The username of the user.",
		},
	}
}

var keysListAccess = []string{"Name"}

type paramsListAccess struct {
	instanceID string
	userName   string
}

type commandListAccess handler.Command

func actionListAccess(c *cli.Context) {
	command := &commandListAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandListAccess) Context() *handler.Context {
	return command.Ctx
}

func (command *commandListAccess) Keys() []string {
	return keysListAccess
}

func (command *commandListAccess) ServiceClientType() string {
	return serviceClientType
}

func (command *commandListAccess) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"instance", "name"})
	if err != nil {
		return err
	}
	c := command.Ctx.CLIContext
	resource.Params = &paramsListAccess{
		instanceID: c.String("instance"),
		userName:   c.String("name"),
	}
	return nil
}

func (command *commandListAccess) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsListAccess)
	pager := users.ListAccess(command.Ctx.ServiceClient, params.instanceID, params.userName)

	pages, err := pager.AllPages()
	if err != nil {
		resource.Err = err
		return
	}

	info, err := users.ExtractDBs(pages)
	if err != nil {
		resource.Err = err
		return
	}

	if len(info) > 0 {
		resource.Result = fmt.Sprintf("%#v", info)
		result := make([]map[string]interface{}, len(info))
		for j, db := range info {
			result[j] = map[string]interface{}{"Name": db.Name}
		}
		resource.Result = result
	} else {
		resource.Result = fmt.Sprintf("User does not has access to any databases")
	}
}
