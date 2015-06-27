package flavorcommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/pagination"
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
		cli.BoolFlag{
			Name:  "all-pages",
			Usage: "[optional] Return all flavors. Default is to paginate.",
		},
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
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Return at most this many flavors.",
		},
	}
}

var keysList = []string{"ID", "Name", "RAM", "Disk", "Swap", "VCPUs", "RxTxFactor"}

type paramsList struct {
	opts     *flavors.ListOpts
	allPages bool
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
		Limit:   c.Int("limit"),
	}
	resource.Params = &paramsList{
		opts:     opts,
		allPages: c.Bool("all-pages"),
	}
	return nil
}

func (command *commandList) HandleSingle(resource *handler.Resource) error {
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	allPages := resource.Params.(*paramsList).allPages
	pager := flavors.ListDetail(command.Ctx.ServiceClient, opts)
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		info, err := flavors.ExtractFlavors(pages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(info))
		for j, flavor := range info {
			result[j] = structs.Map(flavor)
		}
		resource.Result = result
	} else {
		limit := opts.Limit
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			info, err := flavors.ExtractFlavors(page)
			if err != nil {
				return false, err
			}
			result := make([]map[string]interface{}, len(info))
			for j, flavor := range info {
				result[j] = structs.Map(flavor)
			}
			resource.Result = result
			if len(info) >= opts.Limit {
				return false, nil
			}
			limit -= len(info)
			command.Ctx.WaitGroup.Add(1)
			command.Ctx.Results <- resource
			return true, nil
		})
		if err != nil {
			resource.Err = err
			return
		}
	}
}
