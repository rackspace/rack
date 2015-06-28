package securitygrouprulecommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osSecurityGroupRules "github.com/rackspace/gophercloud/openstack/networking/v2/extensions/security/rules"
	"github.com/rackspace/gophercloud/pagination"
	securityGroupRules "github.com/rackspace/gophercloud/rackspace/networking/v2/security/rules"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing security groups",
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
			Usage: "[optional] Return all subnets. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "direction",
			Usage: "[optional] Only list security group rules with this direction. Options are: ingress, egress.",
		},
		cli.StringFlag{
			Name:  "ether-type",
			Usage: "[optional] Only list security group rules with this ether type. Options are: ipv4, ipv6.",
		},
		cli.IntFlag{
			Name:  "port-range-min",
			Usage: "[optional] Only list security group rules that have the low port greater than this.",
		},
		cli.IntFlag{
			Name:  "port-range-max",
			Usage: "[optional] Only list security group rules that have the high port less than this.",
		},
		cli.StringFlag{
			Name:  "protocol",
			Usage: "[optional] Only list security group rules with this protocol. Examples: tcp, udp, icmp.",
		},
		cli.StringFlag{
			Name:  "security-group-id",
			Usage: "[optional] Only list security group rules with this security group ID.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] Only list security group rules that are owned by the tenant with this tenant ID.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing security group rules at this ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many security group rules at most.",
		},
	}
}

var keysList = []string{"ID", "Direction", "EtherType", "PortRangeMin", "PortRangeMax", "Protocol", "SecurityGroupID"}

type paramsList struct {
	opts     *osSecurityGroupRules.ListOpts
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

	opts := &osSecurityGroupRules.ListOpts{
		Direction:    c.String("direction"),
		PortRangeMax: c.Int("port-range-max"),
		PortRangeMin: c.Int("port-range-min"),
		Protocol:     c.String("protocol"),
		SecGroupID:   c.String("security-group-id"),
		TenantID:     c.String("tenant-id"),
		Marker:       c.String("marker"),
		Limit:        c.Int("limit"),
	}

	if c.IsSet("ether-type") {
		etherType := c.String("ether-type")
		switch etherType {
		case "ipv4":
			opts.EtherType = osSecurityGroupRules.Ether4
		case "ipv6":
			opts.EtherType = osSecurityGroupRules.Ether6
		default:
			return fmt.Errorf("Invalid value for `ether-type`: %s. Options are: ipv4, ipv6", etherType)
		}
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
	pager := securityGroupRules.List(command.Ctx.ServiceClient, *opts)
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		info, err := osSecurityGroupRules.ExtractRules(pages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(info))
		for j, rule := range info {
			result[j] = securityGroupRuleSingle(&rule)
		}
		resource.Result = result
	} else {
		limit := opts.Limit
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			info, err := osSecurityGroupRules.ExtractRules(page)
			if err != nil {
				return false, err
			}
			result := make([]map[string]interface{}, len(info))
			for j, rule := range info {
				result[j] = securityGroupRuleSingle(&rule)
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
