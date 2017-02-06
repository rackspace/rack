package securitygrouprulecommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osSecurityGroupRules "github.com/rackspace/gophercloud/openstack/networking/v2/extensions/security/rules"
	securityGroupRules "github.com/rackspace/gophercloud/rackspace/networking/v2/security/rules"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--security-group-id <securityGroupID> --direction <ingress|egress> --ether-type <ipv4|ipv6>"),
	Description: "Creates a security group rule",
	Action:      actionCreate,
	Flags:       commandoptions.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "security-group-id",
			Usage: "[required] The security group ID with which to associate this security group rule.",
		},
		cli.StringFlag{
			Name:  "direction",
			Usage: "[required] The direction of the security group rule. Options are: ingress, egress.",
		},
		cli.StringFlag{
			Name:  "ether-type",
			Usage: "[required] The ether type of the security group rule. Options are: ipv4, ipv6.",
		},
		cli.IntFlag{
			Name:  "port-range-min",
			Usage: "[optional] The minimum port of security group rule.",
		},
		cli.IntFlag{
			Name:  "port-range-max",
			Usage: "[optional] The maximum port of security group rule.",
		},
		cli.StringFlag{
			Name:  "protocol",
			Usage: "[optional] The protocol of the security group rule. Examples: tcp, udp, icmp.",
		},
		cli.StringFlag{
			Name:  "remote-ip-prefix",
			Usage: "[optional] The remote IP prefix to associate with this security group rule",
		},
	}
}

var keysCreate = []string{"ID", "Direction", "EtherType", "PortRangeMin", "PortRangeMax", "Protocol", "SecurityGroupID"}

type paramsCreate struct {
	opts *osSecurityGroupRules.CreateOpts
}

type commandCreate handler.Command

func actionCreate(c *cli.Context) {
	command := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandCreate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandCreate) Keys() []string {
	return keysCreate
}

func (command *commandCreate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"ether-type", "direction", "security-group-id"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext

	opts := &osSecurityGroupRules.CreateOpts{
		PortRangeMax:   c.Int("port-range-max"),
		PortRangeMin:   c.Int("port-range-min"),
		Protocol:       c.String("protocol"),
		SecGroupID:     c.String("security-group-id"),
		RemoteIPPrefix: c.String("remote-ip-prefix"),
	}

	direction := c.String("direction")
	switch direction {
	case "ingress":
		opts.Direction = osSecurityGroupRules.DirIngress
	case "egress":
		opts.Direction = osSecurityGroupRules.DirEgress
	default:
		return fmt.Errorf("Invalid value for `direction`: %s. Options are: ingress, egress", direction)
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

	resource.Params = &paramsCreate{
		opts: opts,
	}

	return nil
}

func (command *commandCreate) HandleSingle(resource *handler.Resource) error {
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	securityGroupRule, err := securityGroupRules.Create(command.Ctx.ServiceClient, *opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = securityGroupRuleSingle(securityGroupRule)
}
