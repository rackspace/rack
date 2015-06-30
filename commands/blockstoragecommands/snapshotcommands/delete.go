package snapshotcommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osSnapshots "github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <snapshotID> | --name <snapshotName> | --stdin id]"),
	Description: "Deletes a snapshot",
	Action:      actionDelete,
	Flags:       util.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the snapshot.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the snapshot.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	snapshotID string
}

type commandDelete handler.Command

func actionDelete(c *cli.Context) {
	command := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandDelete) Context() *handler.Context {
	return command.Ctx
}

func (command *commandDelete) Keys() []string {
	return keysDelete
}

func (command *commandDelete) ServiceClientType() string {
	return serviceClientType
}

func (command *commandDelete) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsDelete{}
	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsDelete).snapshotID = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	snapshotID, err := command.Ctx.IDOrName(osSnapshots.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).snapshotID = snapshotID
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	snapshotID := resource.Params.(*paramsDelete).snapshotID
	err := osSnapshots.Delete(command.Ctx.ServiceClient, snapshotID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Deleting snapshot [%s]\n", snapshotID)
}

func (command *commandDelete) StdinField() string {
	return "id"
}
