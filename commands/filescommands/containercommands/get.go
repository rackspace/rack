package containercommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--name <containerName> | --stdin name]"),
	Description: "Retreives a container",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
	}
}

var keysGet = []string{"Name", "ObjectCount", "BytesUsed", "ContentLength", "AcceptRanges",
	"ContentType", "Read", "Write", "TransID", "VersionsLocation", "Metadata"}

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
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	containerName := command.Ctx.CLIContext.String("name")
	resource.Params.(*paramsGet).container = containerName
	return err
}

func (command *commandGet) Execute(resource *handler.Resource) {
	containerName := resource.Params.(*paramsGet).container
	containerRaw := containers.Get(command.Ctx.ServiceClient, containerName)
	containerInfo, err := containerRaw.Extract()
	if err != nil {
		resource.Err = err
		return
	}
	containerMetadata, err := containerRaw.ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(containerInfo)
	resource.Result.(map[string]interface{})["Name"] = containerName
	resource.Result.(map[string]interface{})["Metadata"] = containerMetadata
}

func (command *commandGet) StdinField() string {
	return "name"
}

func (command *commandGet) PreCSV(resource *handler.Resource) error {
	resource.FlattenMap("Metadata")
	return nil
}

func (command *commandGet) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
