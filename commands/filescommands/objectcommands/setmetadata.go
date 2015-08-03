package objectcommands

import (
	"fmt"

	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	osObjects "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
	"github.com/jrperritt/rack/util"
)

var setMetadata = cli.Command{
	Name:        "set-metadata",
	Usage:       util.Usage(commandPrefix, "set-metadata", "--name <objectName> --container <containerName> --metadata <metadata>"),
	Description: "Set metadata for the given object. This will erase any current metadata.",
	Action:      actionSetMetadata,
	Flags:       commandoptions.CommandFlags(flagsSetMetadata, keysSetMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsSetMetadata, keysSetMetadata))
	},
}

func flagsSetMetadata() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The object name to associate with the metadata.",
		},
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container that holds the object.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[required] A comma-separated string of 'key=value' pairs to set as metadata for the object.",
		},
	}
}

var keysSetMetadata = []string{}

type paramsSetMetadata struct {
	objectName    string
	containerName string
	metadata      map[string]string
}

type commandSetMetadata handler.Command

func actionSetMetadata(c *cli.Context) {
	command := &commandSetMetadata{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandSetMetadata) Context() *handler.Context {
	return command.Ctx
}

func (command *commandSetMetadata) Keys() []string {
	return keysSetMetadata
}

func (command *commandSetMetadata) ServiceClientType() string {
	return serviceClientType
}

func (command *commandSetMetadata) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name", "container", "metadata"})
	if err != nil {
		return err
	}

	metadata, err := command.Ctx.CheckKVFlag("metadata")
	if err != nil {
		return err
	}

	resource.Params = &paramsSetMetadata{
		objectName:    command.Ctx.CLIContext.String("name"),
		containerName: command.Ctx.CLIContext.String("container"),
		metadata:      metadata,
	}
	return err
}

func (command *commandSetMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsSetMetadata)
	containerName := params.containerName
	objectName := params.objectName

	fmt.Printf("params.metadata: %+v\n", params.metadata)

	currentMetadata, err := objects.Get(command.Ctx.ServiceClient, containerName, objectName, nil).ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}
	fmt.Printf("current metadata: %+v\n", currentMetadata)

	for k := range currentMetadata {
		currentMetadata[k] = ""
	}
	fmt.Printf("current metadata emptied: %+v\n", currentMetadata)

	for k, v := range params.metadata {
		currentMetadata[k] = v
	}
	fmt.Printf("current metadata to set: %+v\n", currentMetadata)

	updateOpts := osObjects.UpdateOpts{
		Metadata: currentMetadata,
	}
	updateResponse := objects.Update(command.Ctx.ServiceClient, containerName, objectName, updateOpts)
	if updateResponse.Err != nil {
		resource.Err = updateResponse.Err
		return
	}

	metadata, err := objects.Get(command.Ctx.ServiceClient, containerName, objectName, nil).ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = metadata
}

func (command *commandSetMetadata) PreCSV(resource *handler.Resource) {
	resource.Result = map[string]interface{}{
		"Metadata": resource.Result,
	}
	resource.Keys = []string{"Metadata"}
	resource.FlattenMap("Metadata")
}

func (command *commandSetMetadata) PreTable(resource *handler.Resource) {
	command.PreCSV(resource)
}
