package volumecommands

import (
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osVolumes "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--size <size>"),
	Description: "Creates a volume",
	Action:      actionCreate,
	Flags:       util.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  "size",
			Usage: "[required] The size of this volume (in gigabytes). Valid values are between 75 and 1024.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] A name for this volume.",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] A description for this volume.",
		},
		cli.StringFlag{
			Name:  "volume-type",
			Usage: "[optional] The volume type of this volume.",
		},
	}
}

var keysCreate = []string{"ID", "Name", "Description", "Size", "VolumeType", "SnapshotID", "Attachments", "CreatedAt"}

type paramsCreate struct {
	opts *osVolumes.CreateOpts
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
	err := command.Ctx.CheckFlagsSet([]string{"size"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext

	opts := &osVolumes.CreateOpts{
		Size:        c.Int("size"),
		Name:        c.String("name"),
		Description: c.String("description"),
		VolumeType:  c.String("volume-type"),
	}

	resource.Params = &paramsCreate{
		opts: opts,
	}

	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	volume, err := osVolumes.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = volumeSingle(volume)
}
