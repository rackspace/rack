package usercommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/users"
	"github.com/rackspace/rack/util"
)

var revokeAccess = cli.Command{
	Name:        "revoke-access",
	Usage:       util.Usage(commandPrefix, "revoke-access", "--instance <instanceId> [--name <name> | --stdin name ] --database <database>"),
	Description: "Revokes user access from a database",
	Action:      actionRevokeAccess,
	Flags:       commandoptions.CommandFlags(flagsRevokeAccess, keysRevokeAccess),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsRevokeAccess, keysRevokeAccess))
	},
}

func flagsRevokeAccess() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "instance",
			Usage: "[required] The ID of the instance that the user will be grantAccessd on.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the user.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name.",
		},
		cli.StringFlag{
			Name:  "database",
			Usage: "[required] The database that the user will have access revoked from.",
		},
	}
}

var keysRevokeAccess = []string{}

type commandRevokeAccess handler.Command

func actionRevokeAccess(c *cli.Context) {
	command := &commandRevokeAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandRevokeAccess) Context() *handler.Context {
	return command.Ctx
}

func (command *commandRevokeAccess) Keys() []string {
	return keysRevokeAccess
}

func (command *commandRevokeAccess) ServiceClientType() string {
	return serviceClientType
}

type paramsRevokeAccess struct {
	instanceID string
	userName   string
	database   string
}

func (command *commandRevokeAccess) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	resource.Params = &paramsRevokeAccess{
		instanceID: c.String("instance"),
		userName:   c.String("name"),
		database:   c.String("database"),
	}

	return nil
}

func (command *commandRevokeAccess) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsRevokeAccess).userName = item
	return nil
}

func (command *commandRevokeAccess) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"instance", "name", "database"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsRevokeAccess).userName = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandRevokeAccess) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsRevokeAccess)

	err := users.RevokeAccess(command.Ctx.ServiceClient, params.instanceID, params.userName, params.database).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = fmt.Sprintf("User access has been revoked from: %s", command.Ctx.CLIContext.String("database"))
}
