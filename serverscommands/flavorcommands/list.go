package flavorcommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/flavors"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", util.IDOrNameUsage("flavor")),
	Description: "Lists flavors",
	Action:      actionList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{

		cli.IntFlag{
			Name:  "min-disk",
			Usage: "[optional] Only list flavors that have at least this much disk storage (in GB).",
		},
		cli.IntFlag{
			Name:  "min-ram",
			Usage: "[optional] Only list flavors that have at least this much RAM (in GB).",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing flavors at this flavor ID.",
		},
	}
}

var keysList = []string{"ID", "Name", "RAM", "Disk", "Swap", "VCPUs", "RxTxFactor"}

type paramsList struct {
	opts *flavors.ListOpts
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
	opts := &flavors.ListOpts{
		MinDisk: c.Int("min-disk"),
		MinRAM:  c.Int("min-ram"),
		Marker:  c.String("marker"),
	}
	resource.Params = &paramsList{
		opts: opts,
	}
	return nil
}

func (command *commandList) HandleSingle(resource *handler.Resource) error {
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	allPages, err := flavors.ListDetail(command.Ctx.ServiceClient, opts).AllPages()
	if err != nil {
		resource.Err = err
		return
	}
	flavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		resource.Err = err
		return
	}
	result := make([]map[string]interface{}, len(flavors))
	for j, flavor := range flavors {
		result[j] = structs.Map(flavor)
	}
	resource.Result = result
}
