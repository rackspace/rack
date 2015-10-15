package objectcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "--container <containerName> --name <objectName>"),
	Description: "Retreives an object",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
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

var keysGet = []string{"Name", "ContentDisposition", "ContentEncoding", "ContentLength",
	"ContentType", "StaticLargeObject", "ObjectManifest", "TransID", "Metadata"}

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

	c := command.Ctx.CLIContext
	containerName := c.String("container")
	if err := CheckContainerExists(command.Ctx.ServiceClient, containerName); err != nil {
		return err
	}

	object := c.String("name")
	resource.Params = &paramsGet{
		container: containerName,
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
	objectRaw := objects.Get(command.Ctx.ServiceClient, containerName, objectName, nil)
	objectInfo, err := objectRaw.Extract()
	if err != nil {
		resource.Err = err
		return
	}
	objectMetadata, err := objectRaw.ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(objectInfo)
	resource.Result.(map[string]interface{})["Name"] = objectName
	resource.Result.(map[string]interface{})["Metadata"] = objectMetadata
}

func (command *commandGet) PreCSV(resource *handler.Resource) error {
	resource.FlattenMap("Metadata")
	return nil
}

func (command *commandGet) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
