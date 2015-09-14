package stackcommands

import (
	"errors"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
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
			Name:  "parameter-file",
			Usage: "[optional] Parameter values from file used to create the stack. This can be specified multiple times. Parameter value would be the content of the file",
		},
		cli.StringSliceFlag{
			Name:  "tags",
			Usage: "[optional] A list of tags to associate with the stack.",
		},
	}
}

type paramsUpdate struct {
	opts      *osStacks.UpdateOpts
	stackID   string
	stackName string
}

var keysUpdate = keysList

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
		opts.Tags = c.StringSlice("tags")
	}

	resource.Params = &paramsUpdate{
		opts: opts,
	}
	return nil
}

func (command *commandUpdate) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params.(*paramsUpdate).stackName = name
	resource.Params.(*paramsUpdate).stackID = id
	return nil
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsUpdate)
	opts := params.opts
	stackName := params.stackName
	stackID := params.stackID
	stacks.Update(command.Ctx.ServiceClient, stackName, stackID, opts)
	// the behavior of the python-heatclient is to show a list of stacks as the
	// output of stack-create.
	result, err := stackList(command.Ctx.ServiceClient)
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = result
}
