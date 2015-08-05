package securitygrouprulecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	securityGroupRules "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/networking/v2/security/rules"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", ""),
	Description: "Gets an existing security group rule",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the security group rule.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysGet = []string{"ID", "Direction", "EtherType", "PortRangeMin", "PortRangeMax", "Protocol", "SecurityGroupID", "RemoteGroupID", "RemoteIPPrefix", "TenantID"}

type paramsGet struct {
	securityGroupRuleID string
}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGet).securityGroupRuleID = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).securityGroupRuleID = command.Ctx.CLIContext.String("id")
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	securityGroupRuleID := resource.Params.(*paramsGet).securityGroupRuleID
	securityGroupRule, err := securityGroupRules.Get(command.Ctx.ServiceClient, securityGroupRuleID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = securityGroupRuleSingle(securityGroupRule)
}

func (command *commandGet) StdinField() string {
	return "id"
}
