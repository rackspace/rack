package volumecommands

import (
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osVolumes "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/pagination"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing volumes",
	Action:      actionList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Only list volumes with this name.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Only list volumes that have this status.",
		},
	}
}

var keysList = []string{"ID", "Name", "Size", "Status", "Description", "VolumeType", "SnapshotID", "Attachments", "Created"}

type paramsList struct {
	opts *osVolumes.ListOpts
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

	opts := &osVolumes.ListOpts{
		Name:   c.String("name"),
		Status: c.String("status"),
	}

	resource.Params = &paramsList{
		opts: opts,
	}

	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	pager := osVolumes.List(command.Ctx.ServiceClient, opts)
	var volumes []map[string]interface{}
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		info, err := osVolumes.ExtractVolumes(page)
		if err != nil {
			return false, err
		}
		result := make([]map[string]interface{}, len(info))
		for j, volume := range info {
			result[j] = volumeSingle(&volume)
		}
		volumes = append(volumes, result...)
		return true, nil
	})
	if err != nil {
		resource.Err = err
		return
	}
	if len(volumes) == 0 {
		resource.Result = nil
	} else {
		resource.Result = volumes
	}
}
