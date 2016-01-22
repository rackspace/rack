package instancecommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	"github.com/rackspace/rack/util"
)

var restart = cli.Command{
	Name:        "restart",
	Usage:       util.Usage(commandPrefix, "restart", "[--id <instanceId> | --stdin id]"),
	Description: "Restarts an existing database instance",
	Action:      actionRestart,
	Flags:       commandoptions.CommandFlags(flagsRestart, keysRestart),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsRestart, keysRestart))
	},
}

func flagsRestart() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the instance.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the instance has been restartd",
		},
	}
}

var keysRestart = []string{}

type paramsRestart struct {
	wait bool
	id   string
}

type commandRestart handler.Command

func actionRestart(c *cli.Context) {
	command := &commandRestart{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandRestart) Context() *handler.Context {
	return command.Ctx
}

func (command *commandRestart) Keys() []string {
	return keysRestart
}

func (command *commandRestart) ServiceClientType() string {
	return serviceClientType
}

func (command *commandRestart) HandleFlags(resource *handler.Resource) error {
	wait := false
	if command.Ctx.CLIContext.IsSet("wait-for-completion") {
		wait = true
	}

	resource.Params = &paramsRestart{
		wait: wait,
	}
	return nil
}

func (command *commandRestart) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsRestart).id = item
	return nil
}

func (command *commandRestart) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsRestart).id = command.Ctx.CLIContext.String("id")
	return err
}

func (command *commandRestart) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsRestart)
	err := instances.Restart(command.Ctx.ServiceClient, params.id).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}

	if params.wait {
		err = gophercloud.WaitFor(600, func() (bool, error) {
			inst, err := instances.Get(command.Ctx.ServiceClient, params.id).Extract()
			if err != nil {
				return false, err
			}
			if inst.Status == "ACTIVE" {
				resource.Result = fmt.Sprintln("Restarted instance")
				return true, nil
			}
			return false, nil
		})
	} else {
		resource.Result = fmt.Sprintln("Restarting instance")
	}
}

func (command *commandRestart) StdinField() string {
	return "id"
}
