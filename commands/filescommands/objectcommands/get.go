package objectcommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "--container <containerName> --name <objectName>"),
	Description: "Retreives an object",
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
			Usage: "[required] The name of the container",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the object",
		},
	}
}

var keysGet = []string{"Name", "ContentLength", "ContentType", "StaticLargeObject"}

type paramsGet struct {
	container string
	object    string
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
	err := command.Ctx.CheckFlagsSet([]string{"container", "name"})
	if err != nil {
		return err
	}
	container := command.Ctx.CLIContext.String("container")
	object := command.Ctx.CLIContext.String("name")
	resource.Params = &paramsGet{
		container: container,
		object:    object,
	}
	return nil
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGet).object = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).object = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	containerName := resource.Params.(*paramsGet).container
	objectName := resource.Params.(*paramsGet).object
	objectInfo, err := objects.Get(command.Ctx.ServiceClient, containerName, objectName, nil).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(objectInfo)
	resource.Result.(map[string]interface{})["Name"] = objectName
}

func (command *commandGet) StdinField() string {
	return "name"
}
