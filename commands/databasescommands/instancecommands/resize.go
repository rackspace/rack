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

var resizeFlavor = cli.Command{
	Name:        "resize",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <instanceId> | --stdin id]"),
	Description: "Resizes an existing database instance",
	Action:      actionResize,
	Flags:       commandoptions.CommandFlags(flagsResize, keysResize),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsResize, keysResize))
	},
}

func flagsResize() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the instance.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
		cli.StringFlag{
			Name:  "flavor",
			Usage: "[required] UUID of new flavor",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the instance has been deleted",
		},
	}
}

var keysResize = []string{}

type paramsResize struct {
	wait   bool
	id     string
	flavor string
}

type commandResize handler.Command

func actionResize(c *cli.Context) {
	command := &commandResize{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandResize) Context() *handler.Context {
	return command.Ctx
}

func (command *commandResize) Keys() []string {
	return keysResize
}

func (command *commandResize) ServiceClientType() string {
	return serviceClientType
}

func (command *commandResize) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	wait := false
	if c.IsSet("wait-for-completion") {
		wait = true
	}

	resource.Params = &paramsResize{
		wait:   wait,
		flavor: c.String("flavor"),
	}
	return nil
}

func (command *commandResize) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsResize).id = item
	return nil
}

func (command *commandResize) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id", "flavor"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsResize).id = command.Ctx.CLIContext.String("id")
	return err
}

func (command *commandResize) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsResize)
	id := params.id
	err := instances.Resize(command.Ctx.ServiceClient, id, params.flavor).ExtractErr()
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
				resource.Result = fmt.Sprintf("Resized instance %s", params.id)
				return true, nil
			}
			return false, nil
		})
	} else {
		resource.Result = fmt.Sprintf("Resizing instance %s", params.id)
	}
}

func (command *commandResize) StdinField() string {
	return "id"
}
