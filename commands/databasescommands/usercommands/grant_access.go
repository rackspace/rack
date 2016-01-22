package usercommands

import (
	"fmt"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/users"
	"github.com/rackspace/rack/util"
)

var grantAccess = cli.Command{
	Name:        "grant-access",
	Usage:       util.Usage(commandPrefix, "grant-access", "--instance <instanceId> [--name <name> | --stdin name ] --databases <db1,db2>"),
	Description: "Grants user access to a list of databases",
	Action:      actionGrantAccess,
	Flags:       commandoptions.CommandFlags(flagsGrantAccess, keysGrantAccess),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGrantAccess, keysGrantAccess))
	},
}

func flagsGrantAccess() []cli.Flag {
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
			Name:  "databases",
			Usage: "[optional] A comma-delimeted set of databases that the user will have access granted to.",
		},
	}
}

var keysGrantAccess = []string{}

type commandGrantAccess handler.Command

func actionGrantAccess(c *cli.Context) {
	command := &commandGrantAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGrantAccess) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGrantAccess) Keys() []string {
	return keysGrantAccess
}

func (command *commandGrantAccess) ServiceClientType() string {
	return serviceClientType
}

type paramsGrantAccess struct {
	instanceID string
	userName   string
	opts       *databases.BatchCreateOpts
}

func (command *commandGrantAccess) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := databases.BatchCreateOpts{}

	dbs := strings.Split(c.String("databases"), ",")
	for _, db := range dbs {
		opts = append(opts, databases.CreateOpts{
			Name: db,
		})
	}

	resource.Params = &paramsGrantAccess{
		instanceID: c.String("instance"),
		userName:   c.String("name"),
		opts:       &opts,
	}

	return nil
}

func (command *commandGrantAccess) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGrantAccess).userName = item
	return nil
}

func (command *commandGrantAccess) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"instance", "name", "databases"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGrantAccess).userName = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandGrantAccess) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGrantAccess)

	err := users.GrantAccess(command.Ctx.ServiceClient, params.instanceID, params.userName, *params.opts).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = fmt.Sprintf("User has been granted access to: %s", command.Ctx.CLIContext.String("databases"))
}
