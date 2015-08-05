package objectcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
	"github.com/rackspace/rack/util"
)

var getMetadata = cli.Command{
	Name:        "get-metadata",
	Usage:       util.Usage(commandPrefix, "get-metadata", "--name <objectName> --container <containerName>"),
	Description: "Get metadata for the given object.",
	Action:      actionGetMetadata,
	Flags:       commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata))
	},
}

func flagsGetMetadata() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The object name with the metadata.",
		},
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container that holds the object.",
		},
	}
}

var keysGetMetadata = []string{"Metadata"}

type paramsGetMetadata struct {
	objectName    string
	containerName string
}

type commandGetMetadata handler.Command

func actionGetMetadata(c *cli.Context) {
	command := &commandGetMetadata{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGetMetadata) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGetMetadata) Keys() []string {
	return keysGetMetadata
}

func (command *commandGetMetadata) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGetMetadata) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name", "container"})
	if err != nil {
		return err
	}

	resource.Params = &paramsGetMetadata{
		objectName:    command.Ctx.CLIContext.String("name"),
		containerName: command.Ctx.CLIContext.String("container"),
	}
	return err
}

func (command *commandGetMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGetMetadata)
	containerName := params.containerName
	objectName := params.objectName

	metadata, err := objects.Get(command.Ctx.ServiceClient, containerName, objectName, nil).ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = metadata
}

func (command *commandGetMetadata) PreCSV(resource *handler.Resource) {
	resource.Result = map[string]interface{}{
		"Metadata": resource.Result,
	}
	resource.Keys = []string{"Metadata"}
	resource.FlattenMap("Metadata")
}

func (command *commandGetMetadata) PreTable(resource *handler.Resource) {
	command.PreCSV(resource)
}
