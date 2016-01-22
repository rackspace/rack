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

var resizeDisk = cli.Command{
	Name:        "resize-disk",
	Usage:       util.Usage(commandPrefix, "resizeDisk", "[--id <instanceId> | --stdin id]"),
	Description: "Resizes the disk of an existing database instance",
	Action:      actionResizeDisk,
	Flags:       commandoptions.CommandFlags(flagsResizeDisk, keysResizeDisk),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsResizeDisk, keysResizeDisk))
	},
}

func flagsResizeDisk() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the instance.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
		cli.IntFlag{
			Name:  "size",
			Usage: "[required] The new volume size in GB",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the instance has been resizeDiskd",
		},
	}
}

var keysResizeDisk = []string{}

type paramsResizeDisk struct {
	wait bool
	id   string
	size int
}

type commandResizeDisk handler.Command

func actionResizeDisk(c *cli.Context) {
	command := &commandResizeDisk{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandResizeDisk) Context() *handler.Context {
	return command.Ctx
}

func (command *commandResizeDisk) Keys() []string {
	return keysResizeDisk
}

func (command *commandResizeDisk) ServiceClientType() string {
	return serviceClientType
}

func (command *commandResizeDisk) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	wait := false
	if c.IsSet("wait-for-completion") {
		wait = true
	}

	resource.Params = &paramsResizeDisk{
		wait: wait,
		size: c.Int("size"),
	}
	return nil
}

func (command *commandResizeDisk) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsResizeDisk).id = item
	return nil
}

func (command *commandResizeDisk) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id", "size"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsResizeDisk).id = command.Ctx.CLIContext.String("id")
	return err
}

func (command *commandResizeDisk) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsResizeDisk)
	id := params.id
	err := instances.ResizeVolume(command.Ctx.ServiceClient, id, params.size).ExtractErr()
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

func (command *commandResizeDisk) StdinField() string {
	return "id"
}
