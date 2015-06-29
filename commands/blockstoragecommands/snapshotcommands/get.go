package snapshotcommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osSnapshots "github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--id <snapshotID> | --stdin id]"),
	Description: "Gets a snapshot",
	Action:      actionGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the snapshot.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysGet = []string{"ID", "Name", "Size", "Status", "VolumeID", "VolumeType", "SnapshotID", "Bootable", "Attachments"}

type paramsGet struct {
	snapshotID string
}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGet).snapshotID = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).snapshotID = command.Ctx.CLIContext.String("id")
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	snapshotID := resource.Params.(*paramsGet).snapshotID
	snapshot, err := osSnapshots.Get(command.Ctx.ServiceClient, snapshotID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = snapshotSingle(snapshot)
}

func (command *commandGet) StdinField() string {
	return "id"
}
