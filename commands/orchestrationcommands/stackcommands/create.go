package stackcommands

import (
	"errors"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "[--name <stackName>] [--template-file <templateFile> | --template-url <templateURL>]"),
	Description: "Creates a stack",
	Action:      actionCreate,
	Flags:       commandoptions.CommandFlags(flagsCreate, nil),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
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
			Usage: "[optional] Path to the file or URL containing environment for the stack",
		},
		cli.IntFlag{
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
		cli.StringFlag{
			Name:  "tags",
			Usage: "[optional] A comma-separated string of tags to associate with the stack.",
		},
	}
}

type paramsCreate struct {
	opts *osStacks.CreateOpts
}

var keysCreate = []string{"ID", "Links"}

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
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	c := command.Ctx.CLIContext

	opts := &osStacks.CreateOpts{
		Name:         c.String("name"),
		TemplateOpts: new(osStacks.Template),
	}

	if c.IsSet("disable-rollback") {
		disableRollback := c.BoolT("disable-rollback")
		opts.DisableRollback = &disableRollback
	}

	// check if either template url or template file is set
	if c.IsSet("template-file") {
		opts.TemplateOpts.URL = c.String("template-file")
	} else if c.IsSet("template-url") {
		opts.TemplateOpts.URL = c.String("template-url")
	} else {
		return errors.New("Neither template-file nor template-url specified")
	}

	if c.IsSet("environment-file") {
		opts.EnvironmentOpts = new(osStacks.Environment)
		opts.EnvironmentOpts.URL = c.String("environment-file")
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

	if c.IsSet("tags") {
		opts.Tags = strings.Split(c.String("tags"), ",")
	}

	resource.Params = &paramsCreate{
		opts: opts,
	}
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	stack, err := osStacks.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = stack
}

func (command *commandCreate) PreCSV(resource *handler.Resource) error {
	resource.Result = stackSingle(resource.Result)
	resource.FlattenMap("Links")
	return nil
}

func (command *commandCreate) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
