package objectcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osObjects "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
	"github.com/rackspace/rack/util"
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

	c := command.Ctx.CLIContext
	containerName := c.String("container")
	if err := checkContainerExists(command.Ctx.ServiceClient, containerName); err != nil {
		return err
	}

	metadata, err := command.Ctx.CheckKVFlag("metadata")
	if err != nil {
		return err
	}

	resource.Params = &paramsSetMetadata{
		objectName:    c.String("name"),
		containerName: containerName,
		metadata:      metadata,
	}
	return err
}

func (command *commandSetMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsSetMetadata)
	containerName := params.containerName
	objectName := params.objectName

	currentMetadata, err := objects.Get(command.Ctx.ServiceClient, containerName, objectName, nil).ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}

	for k := range currentMetadata {
		currentMetadata[k] = ""
	}

	for k, v := range params.metadata {
		currentMetadata[k] = v
	}

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

func (command *commandSetMetadata) PreCSV(resource *handler.Resource) error {
	resource.Result = map[string]interface{}{
		"Metadata": resource.Result,
	}
	resource.Keys = []string{"Metadata"}
	resource.FlattenMap("Metadata")
	return nil
}

func (command *commandSetMetadata) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
