package subnetcommands

import (
	"fmt"
	"strconv"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osSubnets "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/networking/v2/subnets"
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
			Usage: "[optional] Return all subnets. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] Only list subnets with this name.",
		},
		cli.StringFlag{
			Name:  "network-id",
			Usage: "[optional] Only list subnets with this network ID.",
		},
		cli.IntFlag{
			Name:  "ip-version",
			Usage: "[optional] Only list subnets that have addresses of this IP version. Options are: 4, 6.",
		},
		cli.StringFlag{
			Name:  "gateway-ip",
			Usage: "[optional] Only list subnets with this gateway IP address.",
		},
		cli.StringFlag{
			Name:  "dhcp-enabled",
			Usage: "[optional] Whether or not to list subnets that are DCP-enabled. Options are: true, false.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] Only list subnets that are owned by the tenant with this tenant ID.",
		},
		cli.StringFlag{
			Name:  "cidr",
			Usage: "[optional] Only list subnets that have this CIDR.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing subnets at this ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many subnets at most.",
		},
	}
}

var keysList = []string{"ID", "Name", "Network ID", "CIDR", "EnableDHCP", "GatewayIP"}

type paramsList struct {
	opts     *osSubnets.ListOpts
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
	opts := &osSubnets.ListOpts{
		Name:      c.String("name"),
		NetworkID: c.String("network-id"),
		IPVersion: c.Int("ip-version"),
		GatewayIP: c.String("gateway-ip"),
		TenantID:  c.String("tenant-id"),
		CIDR:      c.String("cidr"),
		Marker:    c.String("marker"),
		Limit:     c.Int("limit"),
	}
	if c.IsSet("enable-dhcp") {
		dhcpRaw := c.String("enable-dhcp")
		dhcp, err := strconv.ParseBool(dhcpRaw)
		if err != nil {
			return fmt.Errorf("Invalid value for flag `enable-dhcp`: %s. Options are: true, false", dhcpRaw)
		}
		opts.EnableDHCP = &dhcp
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
	pager := subnets.List(command.Ctx.ServiceClient, opts)
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		info, err := osSubnets.ExtractSubnets(pages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(info))
		for j, subnet := range info {
			result[j] = subnetSingle(&subnet)
		}
		resource.Result = result
	} else {
		limit := opts.Limit
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			info, err := osSubnets.ExtractSubnets(page)
			if err != nil {
				return false, err
			}
			result := make([]map[string]interface{}, len(info))
			for j, subnet := range info {
				result[j] = subnetSingle(&subnet)
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
