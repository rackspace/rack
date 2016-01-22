package usercommands

import (
	"fmt"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--instance <instanceId> --name <name>"),
	Description: "Creates a new user on a given instance",
	Action:      actionCreate,
	Flags:       commandoptions.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "instance",
			Usage: "[required] The ID of the instance that the user will be created on.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the user being created.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name.",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "[optional] User password for database access.",
		},
		cli.StringFlag{
			Name:  "databases",
			Usage: "[optional] A comma-delimeted set of databases that the user will have access granted to.",
		},
	}
}

var keysCreate = []string{}

type commandCreate handler.Command

func actionCreate(c *cli.Context) {
	command := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandCreate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandCreate) Keys() []string {
	return keysCreate
}

func (command *commandCreate) ServiceClientType() string {
	return serviceClientType
}

type paramsCreate struct {
	instanceID string
	opts       *users.CreateOpts
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &users.CreateOpts{
		Name:     c.String("name"),
		Password: c.String("password"),
	}

	if c.IsSet("databases") {
		opts.Databases = databases.BatchCreateOpts{}
		dbs := strings.Split(c.String("databases"), ",")
		for _, db := range dbs {
			opts.Databases = append(opts.Databases, databases.CreateOpts{
				Name: db,
			})
		}
	}

	resource.Params = &paramsCreate{
		instanceID: c.String("instance"),
		opts:       opts,
	}

	return nil
}

func (command *commandCreate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsCreate).opts.Name = item
	return nil
}

func (command *commandCreate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"instance", "name", "password"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsCreate).opts.Name = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsCreate)

	opts := &users.BatchCreateOpts{*params.opts}

	err := users.Create(command.Ctx.ServiceClient, params.instanceID, opts).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = fmt.Sprintf("User created")
}
