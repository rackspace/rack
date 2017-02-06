package stackcommands

import (
	"errors"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osStacks "github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", "[--name <stackName> | --id <stackID>] [--template-file <templateFile> | --template-url <templateURL>]"),
	Description: "Updates a specified stack",
	Action:      actionUpdate,
	Flags:       commandoptions.CommandFlags(flagsUpdate, nil),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsUpdate, keysUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` isn't provided] The stack name.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't provided] The stack id.",
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

type paramsUpdate struct {
	opts      *osStacks.UpdateOpts
	stackID   string
	stackName string
}

var keysUpdate = keysGet

type commandUpdate handler.Command

func actionUpdate(c *cli.Context) {
	command := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUpdate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUpdate) Keys() []string {
	return keysUpdate
}

func (command *commandUpdate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUpdate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	opts := &osStacks.UpdateOpts{
		TemplateOpts: new(osStacks.Template),
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

	resource.Params = &paramsUpdate{
		stackName: name,
		stackID:   id,
		opts:      opts,
	}
	return nil
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsUpdate)
	opts := params.opts
	stackName := params.stackName
	stackID := params.stackID
	err := stacks.Update(command.Ctx.ServiceClient, stackName, stackID, opts).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	stack, err := stacks.Get(command.Ctx.ServiceClient, stackName, stackID).Extract()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = stack
}

func (command *commandUpdate) PreCSV(resource *handler.Resource) error {
	resource.Result = stackSingle(resource.Result)
	resource.FlattenMap("Parameters")
	resource.FlattenMap("Outputs")
	resource.FlattenMap("Links")
	resource.FlattenMap("NotificationTopics")
	resource.FlattenMap("Capabilities")
	return nil
}

func (command *commandUpdate) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
