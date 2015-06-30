package snapshotcommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osSnapshots "github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/gophercloud/pagination"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing snapshots",
	Action:      actionList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "volume-id",
			Usage: "Only list snapshots with this volume ID.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "Only list snapshots with this name.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Only list snapshots that have this status.",
		},
	}
}

var keysList = []string{"ID", "Name", "Size", "Status", "VolumeID", "VolumeType", "SnapshotID", "Bootable", "Attachments"}

type paramsList struct {
	opts *osSnapshots.ListOpts
}

type commandList handler.Command

func actionList(c *cli.Context) {
	command := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandList) Context() *handler.Context {
	return command.Ctx
}

func (command *commandList) Keys() []string {
	return keysList
}

func (command *commandList) ServiceClientType() string {
	return serviceClientType
}

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &osSnapshots.ListOpts{
		VolumeID: c.String("volume-id"),
		Name:     c.String("name"),
		Status:   c.String("status"),
	}

	resource.Params = &paramsList{
		opts: opts,
	}

	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	pager := osSnapshots.List(command.Ctx.ServiceClient, opts)
	var snapshots []map[string]interface{}
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		info, err := osSnapshots.ExtractSnapshots(page)
		if err != nil {
			return false, err
		}
		result := make([]map[string]interface{}, len(info))
		for j, snapshot := range info {
			result[j] = snapshotSingle(&snapshot)
		}
		snapshots = append(snapshots, result...)
		return true, nil
	})
	if err != nil {
		resource.Err = err
		return
	}
	if len(snapshots) == 0 {
		resource.Result = nil
	} else {
		resource.Result = snapshots
	}
}
