package stackcommands

import (
	"errors"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

var preview = cli.Command{
	Name:        "preview",
	Usage:       util.Usage(commandPrefix, "preview", "--name stackName [--template-file <templateFile> | --template-url <templateURL>]"),
	Description: "Preview a stack",
	Action:      actionPreview,
	Flags:       commandoptions.CommandFlags(flagsPreview, nil),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsPreview, keysPreview))
	},
}

func flagsPreview() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The stack name.",
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

type paramsPreview struct {
	opts *osStacks.PreviewOpts
}

var keysPreview = []string{"Capabilities", "CreationTime", "Description", "DisableRollback", "ID", "Links", "NotificiationTopics", "Parameters", "Resources", "Name", "TemplateDescription", "Timeout", "UpdatedTime"}

type commandPreview handler.Command

func actionPreview(c *cli.Context) {
	command := &commandPreview{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandPreview) Context() *handler.Context {
	return command.Ctx
}

func (command *commandPreview) Keys() []string {
	return keysPreview
}

func (command *commandPreview) ServiceClientType() string {
	return serviceClientType
}

func (command *commandPreview) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &osStacks.PreviewOpts{
		Name:        c.String("name"),
		TemplateOpts: new(osStacks.Template),
	}

	if c.IsSet("disable-rollback") {
		disableRollback := c.Bool("disable-rollback")
		opts.DisableRollback = &disableRollback
	}

	// check if either template url or template file is set
	if c.IsSet("template-file") {
		opts.TemplateOpts.TemplateURL = c.String("template-file")
	} else if c.IsSet("template-url") {
		opts.TemplateOpts.TemplateURL = c.String("template-url")
	} else {
		return errors.New("Neither template-file nor template-url specified")
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

	resource.Params = &paramsPreview{
		opts: opts,
	}
	return nil
}

func (command *commandPreview) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsPreview).opts.Name = item
	return nil
}

func (command *commandPreview) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsPreview).opts.Name = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandPreview) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsPreview).opts
	stack, err := osStacks.Preview(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = stackSingle(stack)
}
