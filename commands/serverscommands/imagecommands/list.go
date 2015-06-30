package imagecommands

import (
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osImages "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/pagination"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/images"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists images",
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
			Usage: "[optional] Return all images. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] Only list images that have this name.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "[optional] Only list images that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing images at this image ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many images at most.",
		},
	}
}

var keysList = []string{"ID", "Name", "Status", "MinDisk", "MinRAM"}

type paramsList struct {
	opts     *osImages.ListOpts
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
	opts := &osImages.ListOpts{
		Name:   c.String("name"),
		Status: c.String("status"),
		Marker: c.String("marker"),
		Limit:  c.Int("limit"),
	}
	resource.Params = &paramsList{
		opts:     opts,
		allPages: c.Bool("all-pages"),
	}
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	allPages := resource.Params.(*paramsList).allPages
	pager := images.ListDetail(command.Ctx.ServiceClient, opts)
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		info, err := osImages.ExtractImages(pages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(info))
		for j, image := range info {
			result[j] = structs.Map(image)
		}
		resource.Result = result
	} else {
		limit := opts.Limit
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			info, err := osImages.ExtractImages(page)
			if err != nil {
				return false, err
			}
			result := make([]map[string]interface{}, len(info))
			for j, image := range info {
				result[j] = structs.Map(image)
			}
			resource.Result = result
			if len(info) >= limit {
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
