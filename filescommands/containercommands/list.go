package containercommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osContainers "github.com/rackspace/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists containers",
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
			Usage: "[optional] Return all containers. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "prefix",
			Usage: "[optional] Only return containers with this prefix.",
		},
		cli.StringFlag{
			Name:  "end-marker",
			Usage: "[optional] Only return containers with name less than this value.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing containers at this container name.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many containers at most.",
		},
	}
}

var keysList = []string{"Name", "Count", "Bytes"}

type paramsList struct {
	opts     *osContainers.ListOpts
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
	opts := &osContainers.ListOpts{
		Full:      c.Bool("full"),
		Prefix:    c.String("prefix"),
		EndMarker: c.String("end-marker"),
		Marker:    c.String("marker"),
		Limit:     c.Int("limit"),
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
	opts.Full = true
	allPages := resource.Params.(*paramsList).allPages
	pager := containers.List(command.Ctx.ServiceClient, opts)
	var containerInfo []osContainers.Container
	var err error
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		containerInfo, err = containers.ExtractInfo(pages)
	} else {
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			containerInfo, err = containers.ExtractInfo(page)
			if err != nil {
				return false, err
			}
			return false, nil
		})
		if err != nil {
			resource.Err = err
			return
		}
	}
	result := make([]map[string]interface{}, len(containerInfo))
	for j, container := range containerInfo {
		result[j] = structs.Map(&container)
	}
	resource.Result = result
}
