package containercommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
)

var create = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--container <containerName> | --stdin container]"),
	Description: "Retreives a container",
	Action:      actionGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `container` isn't provided] The field being piped into STDIN. Valid values are: container",
		},
	}
}

var keysGet = []string{"Name", "ObjectCount", "BytesUsed", "ContentLength"}

type paramsGet struct {
	container string
}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGet).container = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"container"})
	if err != nil {
		return err
	}
	containerName := command.Ctx.CLIContext.String("container")
	resource.Params = containerName
	return err
}

func (command *commandGet) Execute(resource *handler.Resource) {
	containerName := resource.Params.(*paramsGet).container
	containerInfo, err := containers.Get(command.Ctx.ServiceClient, containerName).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(containerInfo)
	resource.Result.(map[string]interface{})["Name"] = containerName
}

func (command *commandGet) StdinField() string {
	return "container"
}
