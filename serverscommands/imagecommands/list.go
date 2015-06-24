package imagecommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osImages "github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/images"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", util.IDOrNameUsage("image")),
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
			Usage: "Only list images that have this name.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Only list images that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "Start listing images at this image ID.",
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
	pager := images.ListDetail(command.Ctx.ServiceClient, opts)
	var imageInfo []osImages.Image
	var err error
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		imageInfo, err = osImages.ExtractImages(pages)
	} else {
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			imageInfo, err = osImages.ExtractImages(page)
			if err != nil {
				return false, err
			}
			return true, nil
		})
		if err != nil {
			resource.Err = err
			return
		}
	}
	result := make([]map[string]interface{}, len(imageInfo))
	for j, image := range imageInfo {
		result[j] = structs.Map(image)
	}
	resource.Result = result
}
