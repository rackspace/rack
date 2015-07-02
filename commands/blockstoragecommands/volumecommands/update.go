package volumecommands

import (
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	osVolumes "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/jrperritt/rack/util"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", "--id <volumeID>"),
	Description: "Updates a volume",
	Action:      actionUpdate,
	Flags:       util.CommandFlags(flagsUpdate, keysUpdate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsUpdate, keysUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[required] The ID of the volume to update.",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] A new name for this volume.",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] A new description for this volume.",
		},
	}
}

var keysUpdate = []string{"ID", "Name", "Description", "Size", "VolumeType", "SnapshotID", "Attachments"}

type paramsUpdate struct {
	volumeID string
	opts     *osVolumes.UpdateOpts
}

type commandUpdate handler.Command

func actionUpdate(c *cli.Context) {
	command := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUpdate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUpdate) Keys() []string {
	return keysUpdate
}

func (command *commandUpdate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUpdate) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"volume-id"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext

	opts := &osVolumes.UpdateOpts{
		Name:        c.String("rename"),
		Description: c.String("description"),
	}

	resource.Params = &paramsUpdate{
		volumeID: c.String("volume-id"),
		opts:     opts,
	}

	return nil
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsUpdate).opts
	volumeID := resource.Params.(*paramsUpdate).volumeID
	volume, err := osVolumes.Update(command.Ctx.ServiceClient, volumeID, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = volumeSingle(volume)
}
