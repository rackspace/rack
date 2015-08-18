package stackcommands

import (
	"io/ioutil"
	"path/filepath"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

var adopt = cli.Command{
	Name:        "adopt",
	Usage:       util.Usage(commandPrefix, "adopt", "--name stackName --adopt-file adoptFile [--template-file <templateFile> | --template-url <templateURL>]"),
	Description: "Creates a stack from existing resources",
	Action:      actionAdopt,
	Flags:       commandoptions.CommandFlags(flagsAdopt, nil),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsAdopt, keysAdopt))
	},
}

func flagsAdopt() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The stack name.",
		},
		cli.StringFlag{
			Name:  "adopt-file",
			Usage: "[required] Path to adopt stack data file",
		},
		cli.StringFlag{
			Name:  "template-file",
			Usage: "[optional; required if `template-url` isn't provided] The path to template file.",
		},
		cli.StringFlag{
			Name:  "template-url",
			Usage: "[optional; required if `template-file` isn't provided] The url to the template.",
		},
		cli.StringFlag{
			Name:  "environment-file",
			Usage: "[optional] File containing environment for the stack",
		},
		cli.StringFlag{
			Name:  "timeout",
			Usage: "[optional] Stack creation timeout in minutes.",
		},
		cli.BoolTFlag{
			Name:  "disable-rollback",
			Usage: "[optional] Disable rollback on create/update failure.",
		},
		cli.StringFlag{
			Name:  "parameters",
			Usage: "[optional] A comma-separated string of key=value pairs.",
		},
	}
}

type paramsAdopt struct {
	opts *osStacks.AdoptOpts
}

var keysAdopt = []string{"Capabilities", "CreationTime", "Description", "DisableRollback", "ID", "Links", "NotificiationTopics", "Outputs", "Parameters", "Name", "Status", "StatusReason", "TemplateDescription", "Timeout", "UpdatedTime"}

type commandAdopt handler.Command

func actionAdopt(c *cli.Context) {
	command := &commandAdopt{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandAdopt) Context() *handler.Context {
	return command.Ctx
}

func (command *commandAdopt) Keys() []string {
	return keysAdopt
}

func (command *commandAdopt) ServiceClientType() string {
	return serviceClientType
}

func (command *commandAdopt) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &osStacks.AdoptOpts{
		Name:        c.String("name"),
	}

	if c.IsSet("disable-rollback") {
		disableRollback := c.Bool("disable-rollback")
		opts.DisableRollback = &disableRollback
	}

	if c.IsSet("environment-file") {
		opts.EnvironmentOpts = new(osStacks.Environment)
		opts.EnvironmentOpts.EnvironmentURL = c.String("environment-file")
	}

	if c.IsSet("timeout") {
		opts.Timeout = c.Int("timeout")
	}

	if c.IsSet("parameters") {
		parameters, err := command.Ctx.CheckKVFlag("parameters")
		if err != nil {
			return err
		}
		opts.Parameters = parameters
	}

	if c.IsSet("adopt-file") {
		abs, err := filepath.Abs(c.String("adopt-file"))
		if err != nil {
			return err
		}
		environment, err := ioutil.ReadFile(abs)
		if err != nil {
			return err
		}
		opts.AdoptStackData = string(environment)
	}

	resource.Params = &paramsAdopt{
		opts: opts,
	}
	return nil
}

func (command *commandAdopt) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsAdopt).opts.Name = item
	return nil
}

func (command *commandAdopt) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name", "adopt-file"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsAdopt).opts.Name = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandAdopt) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsAdopt).opts
	stack, err := stacks.Adopt(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(stack)
}
