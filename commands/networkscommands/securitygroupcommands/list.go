package securitygroupcommands

import (
	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	osSecurityGroups "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/pagination"
	securityGroups "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/networking/v2/security/groups"
	"github.com/jrperritt/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing security groups",
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
			Usage: "[optional] Return all security groups. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] Only list security groups with this name.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] Only list security groups that are owned by the tenant with this tenant ID.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing security groups at this ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many security groups at most.",
		},
	}
}

var keysList = []string{"ID", "Name", "TenantID"}

type paramsList struct {
	opts     *osSecurityGroups.ListOpts
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

	opts := &osSecurityGroups.ListOpts{
		Name:     c.String("name"),
		TenantID: c.String("tenant-id"),
		Marker:   c.String("marker"),
		Limit:    c.Int("limit"),
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
	pager := securityGroups.List(command.Ctx.ServiceClient, *opts)
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		info, err := osSecurityGroups.ExtractGroups(pages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(info))
		for j, subnet := range info {
			result[j] = securityGroupSingle(&subnet)
		}
		resource.Result = result
	} else {
		limit := opts.Limit
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			info, err := osSecurityGroups.ExtractGroups(page)
			if err != nil {
				return false, err
			}
			result := make([]map[string]interface{}, len(info))
			for j, subnet := range info {
				result[j] = securityGroupSingle(&subnet)
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
