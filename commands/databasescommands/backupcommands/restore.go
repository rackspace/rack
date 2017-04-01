package backupcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	"github.com/rackspace/rack/util"
)

var restore = cli.Command{
	Name:        "restore",
	Usage:       util.Usage(commandPrefix, "create", "[--id <id> | --stdin id] --name <instanceName> --flavor <flavorId> --size <instanceSize> --wait-for-completion"),
	Description: "Restores a backed up database instance",
	Action:      actionRestore,
	Flags:       commandoptions.CommandFlags(flagsRestore, keysRestore),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsRestore, keysRestore))
	},
}

func flagsRestore() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the backup being restored.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name that the database instance should have.",
		},
		cli.StringFlag{
			Name:  "flavor",
			Usage: "[required] The flavor ID that the database instance will be based on",
		},
		cli.IntFlag{
			Name:  "size",
			Usage: "[required] The disk space that will be allocated for the database instance in GB",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the instance is available.",
		},
	}
}

var keysRestore = []string{"ID", "Hostname"}

type commandRestore handler.Command

func actionRestore(c *cli.Context) {
	command := &commandRestore{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandRestore) Context() *handler.Context {
	return command.Ctx
}

func (command *commandRestore) Keys() []string {
	return keysRestore
}

func (command *commandRestore) ServiceClientType() string {
	return serviceClientType
}

type paramsRestore struct {
	wait bool
	opts *instances.CreateOpts
}

func (command *commandRestore) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	resource.Params = &paramsRestore{
		wait: c.IsSet("wait-for-completion"),
		opts: &instances.CreateOpts{
			FlavorRef:    c.String("flavor"),
			Name:         c.String("name"),
			Size:         c.Int("size"),
			RestorePoint: c.String("id"),
		},
	}
	return nil
}

func (command *commandRestore) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsRestore).opts.RestorePoint = item
	return nil
}

func (command *commandRestore) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id", "name", "flavor", "size"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsRestore).opts.RestorePoint = command.Ctx.CLIContext.String("id")
	return nil
}

func (command *commandRestore) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsRestore).opts

	instance, err := instances.Create(command.Ctx.ServiceClient, opts).Extract()

	if resource.Params.(*paramsRestore).wait {
		err = gophercloud.WaitFor(600, func() (bool, error) {
			inst, err := instances.Get(command.Ctx.ServiceClient, instance.ID).Extract()
			if err != nil {
				return false, err
			}
			if inst.Status == "ACTIVE" {
				return true, nil
			}
			return false, nil
		})
	}

	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = singleInstance(&instance)
}
