package containercommands

import (
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	"github.com/rackspace/rack/util"
)

var updateMetadata = cli.Command{
	Name:        "update-metadata",
	Usage:       util.Usage(commandPrefix, "update-metadata", "--name <containerName> --metadata <metadata>"),
	Description: "Create or replace metadata for the given container. Any existing metadata will remain in tact.",
	Action:      actionUpdateMetadata,
	Flags:       commandoptions.CommandFlags(flagsUpdateMetadata, keysUpdateMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsUpdateMetadata, keysUpdateMetadata))
	},
}

func flagsUpdateMetadata() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The container name with the metadata.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[required] A comma-separated string of 'key=value' pairs to create of update as metadata for the container.",
		},
	}
}

var keysUpdateMetadata = []string{}

type paramsUpdateMetadata struct {
	containerName string
	metadata      map[string]string
}

type commandUpdateMetadata handler.Command

func actionUpdateMetadata(c *cli.Context) {
	command := &commandUpdateMetadata{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUpdateMetadata) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUpdateMetadata) Keys() []string {
	return keysUpdateMetadata
}

func (command *commandUpdateMetadata) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUpdateMetadata) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name", "metadata"})
	if err != nil {
		return err
	}

	metadata, err := command.Ctx.CheckKVFlag("metadata")
	if err != nil {
		return err
	}

	resource.Params = &paramsUpdateMetadata{
		containerName: command.Ctx.CLIContext.String("name"),
		metadata:      metadata,
	}
	return err
}

func (command *commandUpdateMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsUpdateMetadata)
	containerName := params.containerName

	getResponse := containers.Get(command.Ctx.ServiceClient, containerName)
	if getResponse.Err != nil {
		resource.Err = getResponse.Err
		return
	}

	updateOpts := containers.UpdateOpts{
		Metadata: params.metadata,
	}
	updateResponse := containers.Update(command.Ctx.ServiceClient, containerName, updateOpts)
	if updateResponse.Err != nil {
		resource.Err = updateResponse.Err
		return
	}

	metadata, err := containers.Get(command.Ctx.ServiceClient, containerName).ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}

	updatedMetadata := make(map[string]string, len(params.metadata))
	for k := range params.metadata {
		k = strings.Title(k)
		updatedMetadata[k] = metadata[k]
	}

	resource.Result = updatedMetadata
}

func (command *commandUpdateMetadata) PreCSV(resource *handler.Resource) {
	resource.Result = map[string]interface{}{
		"Metadata": resource.Result,
	}
	resource.Keys = []string{"Metadata"}
	resource.FlattenMap("Metadata")
}

func (command *commandUpdateMetadata) PreTable(resource *handler.Resource) {
	command.PreCSV(resource)
}
