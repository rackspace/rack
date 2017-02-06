package snapshotcommands

import (
	"fmt"
	"time"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osSnapshots "github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/rack/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <snapshotID> | --name <snapshotName> | --stdin id]"),
	Description: "Deletes a snapshot",
	Action:      actionDelete,
	Flags:       commandoptions.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDelete, keysDelete))
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
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the snapshot is available.",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	wait       bool
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
	wait := false
	if command.Ctx.CLIContext.IsSet("wait-for-completion") {
		wait = true
	}

	resource.Params = &paramsDelete{
		wait: wait,
	}
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

	if resource.Params.(*paramsDelete).wait {
		i := 0
		for i < 120 {
			_, err := osSnapshots.Get(command.Ctx.ServiceClient, snapshotID).Extract()
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
			i++
		}
		resource.Result = fmt.Sprintf("Deleted snapshot [%s]\n", snapshotID)
	} else {
		resource.Result = fmt.Sprintf("Deleting snapshot [%s]\n", snapshotID)
	}
}

func (command *commandDelete) StdinField() string {
	return "id"
}
