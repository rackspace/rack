package snapshotcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osSnapshots "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/blockstorage/v1/snapshots"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--volume-id <volumeID>"),
	Description: "Creates a volume",
	Action:      actionCreate,
	Flags:       commandoptions.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "volume-id",
			Usage: "[required] The volume ID from which to create this snapshot.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] A name for this snapshot.",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] A description for this snapshot.",
		},
		cli.BoolFlag{
			Name:  "force",
			Usage: "[optional] If provided, the command will force the request for volumes in-use. Not recommended, risk of data inconsistency.",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the snapshot is available.",
		},
	}
}

var keysCreate = []string{"ID", "Name", "Description", "Size", "VolumeID", "VolumeType", "SnapshotID", "Attachments", "CreatedAt", "Metadata"}

type paramsCreate struct {
	wait bool
	force bool
	opts *snapshots.CreateOpts
}

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

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"volume-id"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext
	wait := false
	if c.IsSet("wait-for-completion") {
		wait = true
	}

	force := false
	if c.IsSet("force") {
		force = true
	}

	opts := &snapshots.CreateOpts{
		VolumeID:    c.String("volume-id"),
		Name:        c.String("name"),
		Description: c.String("description"),
		Force:       force,
	}

	resource.Params = &paramsCreate{
		wait: wait,
		opts: opts,
	}

	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	snapshot, err := osSnapshots.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}

	if resource.Params.(*paramsCreate).wait {
		err = osSnapshots.WaitForStatus(command.Ctx.ServiceClient, snapshot.ID, "available", 600)
		if err != nil {
			resource.Err = err
			return
		}

		snapshot, err = osSnapshots.Get(command.Ctx.ServiceClient, snapshot.ID).Extract()
		if err != nil {
			resource.Err = err
			return
		}
	}

	resource.Result = snapshotSingle(snapshot)
}

func (command *commandCreate) PreCSV(resource *handler.Resource) {
	resource.FlattenMap("Metadata")
}

func (command *commandCreate) PreTable(resource *handler.Resource) {
	command.PreCSV(resource)
}
