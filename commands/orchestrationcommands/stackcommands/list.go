package stackcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osStacks "github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Retrieve a list of deployed stacks",
	Action:      actionList,
	Flags:       commandoptions.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all-pages",
			Usage: "[optional] Return all stacks. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] Only list stacks with this name.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "[optional] Only list stacks that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing stacks at this stack ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many stacks at most.",
		},
		cli.StringFlag{
			Name:  "sort-key",
			Usage: "[optional] Key used to sort the list of stacks. Valid values are: name, status, created_at, updated_at",
		},
		cli.StringFlag{
			Name:  "sort-dir",
			Usage: "[optional] Specify direction for sort. Valid values are: asc, desc",
		},
	}
}

var keysList = []string{"ID", "Name", "Status", "CreationTime"}

type paramsList struct {
	opts     *osStacks.ListOpts
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
	opts := &osStacks.ListOpts{
		Name:    c.String("name"),
		Status:  c.String("status"),
		Marker:  c.String("marker"),
		Limit:   c.Int("limit"),
		SortKey: osStacks.SortKey(c.String("sort-key")),
		SortDir: osStacks.SortDir(c.String("sort-dir")),
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
	pager := stacks.List(command.Ctx.ServiceClient, opts)
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		info, err := osStacks.ExtractStacks(pages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(info))
		for j, stack := range info {
			result[j] = stackSingle(&stack)
		}
		resource.Result = result
	} else {
		limit := opts.Limit
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			info, err := osStacks.ExtractStacks(page)
			if err != nil {
				return false, err
			}
			result := make([]map[string]interface{}, len(info))
			for j, stack := range info {
				result[j] = stackSingle(&stack)
			}
			resource.Result = result
			if len(info) >= limit {
				return false, nil
			}
			limit -= len(info)
			command.Ctx.Results <- resource
			return true, nil
		})
		if err != nil {
			resource.Err = err
			return
		}
	}
}
