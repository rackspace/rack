package backupcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/backups"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--instance <instanceId> --name <name> --description <description>"),
	Description: "Creates a new backup for a database instance.",
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
			Usage: "[required] The ID of the instance that will be backed up",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the backup.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name.",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] Specifies a long description of the backup.",
		},
	}
}

var keysCreate = []string{"ID", "Name", "Description", "InstanceID", "Created",
	"Updated", "Size", "Status"}

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
	opts       *backups.CreateOpts
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &backups.CreateOpts{
		Name:       c.String("name"),
		InstanceID: c.String("instance"),
	}

	if c.IsSet("description") {
		opts.Description = c.String("description")
	}

	resource.Params = &paramsCreate{opts: opts}
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

	backup, err := backups.Create(command.Ctx.ServiceClient, params.opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = singleBackup(backup)
}
