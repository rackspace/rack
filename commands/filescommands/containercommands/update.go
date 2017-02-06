package containercommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	"github.com/rackspace/rack/util"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", "--name <containerName>"),
	Description: "Updates a container",
	Action:      actionUpdate,
	Flags:       commandoptions.CommandFlags(flagsUpdate, keysUpdate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsUpdate, keysUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the container",
		},
		cli.StringFlag{
			Name:  "container-read",
			Usage: "[optional] Comma-separated list of users for whom to grant read access to the container",
		},
		cli.StringFlag{
			Name:  "container-write",
			Usage: "[optional] Comma-separated list of users for whom to grant write access to the container",
		},
	}
}

var keysUpdate = []string{}

type paramsUpdate struct {
	container string
	opts      containers.UpdateOpts
}

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

	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}

	opts := containers.UpdateOpts{
		ContainerRead:  c.String("container-read"),
		ContainerWrite: c.String("container-write"),
	}
	resource.Params = &paramsUpdate{
		container: c.String("name"),
		opts:      opts,
	}
	return nil
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsUpdate)
	containerName := params.container

	containerRaw := containers.Get(command.Ctx.ServiceClient, containerName)
	if containerRaw.Err != nil {
		resource.Err = containerRaw.Err
		return
	}

	opts := params.opts
	rawResponse := containers.Update(command.Ctx.ServiceClient, containerName, opts)
	if rawResponse.Err != nil {
		resource.Err = rawResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully updated container [%s]\n", containerName)
}
