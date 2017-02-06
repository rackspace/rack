package networkcommands

import (
	"fmt"
	"strconv"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osNetworks "github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/networking/v2/networks"
	"github.com/rackspace/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing networks",
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
			Usage: "[optional] Return all networks. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "Only list networks with this name.",
		},
		cli.StringFlag{
			Name:  "up",
			Usage: "Only list networks that are up or not. Options are: true, false.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "Only list networks that are owned by the tenant with this tenant ID.",
		},
		cli.StringFlag{
			Name:  "shared",
			Usage: "Only list networks that are shared or not. Options are: true, false.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Only list networks that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "Start listing networks at this ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many networks at most.",
		},
	}
}

var keysList = []string{"ID", "Name", "Up", "Status", "Shared", "TenantID"}

type paramsList struct {
	opts     *osNetworks.ListOpts
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
	opts := &osNetworks.ListOpts{
		Name:     c.String("name"),
		TenantID: c.String("tenant-id"),
		Status:   c.String("status"),
		Marker:   c.String("marker"),
		Limit:    c.Int("limit"),
	}
	if c.IsSet("up") {
		upRaw := c.String("up")
		up, err := strconv.ParseBool(upRaw)
		if err != nil {
			return fmt.Errorf("Invalid value for flag `up`: %s. Options are: true, false", upRaw)
		}
		opts.AdminStateUp = &up
	}
	if c.IsSet("shared") {
		sharedRaw := c.String("shared")
		shared, err := strconv.ParseBool(sharedRaw)
		if err != nil {
			return fmt.Errorf("Invalid value for flag `shared`: %s. Options are: true, false", sharedRaw)
		}
		opts.Shared = &shared
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
	pager := networks.List(command.Ctx.ServiceClient, opts)
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		info, err := osNetworks.ExtractNetworks(pages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(info))
		for j, network := range info {
			result[j] = networkSingle(&network)
		}
		resource.Result = result
	} else {
		limit := opts.Limit
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			info, err := osNetworks.ExtractNetworks(page)
			if err != nil {
				return false, err
			}
			result := make([]map[string]interface{}, len(info))
			for j, network := range info {
				result[j] = networkSingle(&network)
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
