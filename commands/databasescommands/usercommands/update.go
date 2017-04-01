package usercommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/users"
	"github.com/rackspace/rack/util"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", "--instance <instanceId> --name <name>"),
	Description: "Updates a new user on a given instance",
	Action:      actionUpdate,
	Flags:       commandoptions.CommandFlags(flagsUpdate, keysUpdate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsUpdate, keysUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "instance",
			Usage: "[required] The ID of the instance.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the existing user.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name.",
		},
		cli.StringFlag{
			Name:  "new-name",
			Usage: "[optional] A new username for the user",
		},
		cli.StringFlag{
			Name:  "new-password",
			Usage: "[optional] The new password for the user.",
		},
		cli.StringFlag{
			Name: "new-host",
			Usage: "[optional] Specifies the host from which a user is allowed to connect to the database. " +
				"Possible values are a string containing an IPv4 address or `%` to allow connecting from any host." +
				"If host is not specified, it defaults to `%`.",
		},
	}
}

var keysUpdate = []string{}

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

type paramsUpdate struct {
	instanceID string
	userName   string
	opts       *users.UpdateOpts
}

func (command *commandUpdate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &users.UpdateOpts{
		Name: c.String("new-name"),
	}

	if c.IsSet("new-password") {
		opts.Password = c.String("new-password")
	}

	if c.IsSet("new-host") {
		opts.Host = c.String("new-host")
	}

	resource.Params = &paramsUpdate{
		instanceID: c.String("instance"),
		userName:   c.String("name"),
		opts:       opts,
	}

	return nil
}

func (command *commandUpdate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsUpdate).opts.Name = item
	return nil
}

func (command *commandUpdate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"instance", "name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsUpdate).opts.Name = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsUpdate)

	err := users.Update(command.Ctx.ServiceClient, params.instanceID, params.userName, *params.opts).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = fmt.Sprintf("User updated")
}
