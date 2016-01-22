package databasecommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--instance <instanceId> --name <name>"),
	Description: "Creates a new database on a given instance",
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
			Usage: "[required] The ID of the instance that the database will be created on.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the database being deleted.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name.",
		},
		cli.IntFlag{
			Name:  "collate",
			Usage: "[optional] Set of rules for comparing characters in a character set. The default value for collate is utf8_general_ci.",
		},
		cli.StringFlag{
			Name:  "charset",
			Usage: "[optional] Set of symbols and encodings. The default character set is utf8.",
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
	opts       *databases.CreateOpts
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &databases.CreateOpts{
		Name: c.String("name"),
	}

	if c.IsSet("charset") {
		opts.CharSet = c.String("charset")
	}

	if c.IsSet("collate") {
		opts.Collate = c.String("collate")
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
	err := command.Ctx.CheckFlagsSet([]string{"name", "instance"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsCreate).opts.Name = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsCreate)

	opts := &databases.BatchCreateOpts{*params.opts}

	err := databases.Create(command.Ctx.ServiceClient, params.instanceID, opts).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = fmt.Sprintf("Database created")
}
