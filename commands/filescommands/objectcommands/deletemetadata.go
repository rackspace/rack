package objectcommands

import (
	"fmt"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osObjects "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
	"github.com/rackspace/rack/util"
)

var deleteMetadata = cli.Command{
	Name:        "delete-metadata",
	Usage:       util.Usage(commandPrefix, "delete-metadata", "--name <objectName> --container <containerName> --metadata-keys <metadataKeys>"),
	Description: "Delete specific metadata from the given object.",
	Action:      actionDeleteMetadata,
	Flags:       commandoptions.CommandFlags(flagsDeleteMetadata, keysDeleteMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDeleteMetadata, keysDeleteMetadata))
	},
}

func flagsDeleteMetadata() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The object name from which to delete the metadata.",
		},
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container that holds the object.",
		},
		cli.StringFlag{
			Name:  "metadata-keys",
			Usage: "[required] A comma-separated string of metadata keys to delete from the container.",
		},
	}
}

var keysDeleteMetadata = []string{}

type paramsDeleteMetadata struct {
	objectName    string
	containerName string
	metadataKeys  []string
}

type commandDeleteMetadata handler.Command

func actionDeleteMetadata(c *cli.Context) {
	command := &commandDeleteMetadata{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandDeleteMetadata) Context() *handler.Context {
	return command.Ctx
}

func (command *commandDeleteMetadata) Keys() []string {
	return keysDeleteMetadata
}

func (command *commandDeleteMetadata) ServiceClientType() string {
	return serviceClientType
}

func (command *commandDeleteMetadata) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name", "container", "metadata-keys"})
	if err != nil {
		return err
	}
	metadataKeys := strings.Split(command.Ctx.CLIContext.String("metadata-keys"), ",")
	for i, k := range metadataKeys {
		metadataKeys[i] = strings.Title(k)
	}

	resource.Params = &paramsDeleteMetadata{
		objectName:    command.Ctx.CLIContext.String("name"),
		containerName: command.Ctx.CLIContext.String("container"),
		metadataKeys:  metadataKeys,
	}
	return nil
}

func (command *commandDeleteMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDeleteMetadata)
	containerName := params.containerName
	objectName := params.objectName

	currentMetadata, err := objects.Get(command.Ctx.ServiceClient, containerName, objectName, nil).ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}

	for _, k := range params.metadataKeys {
		currentMetadata[k] = ""
	}

	updateOpts := osObjects.UpdateOpts{
		Metadata: currentMetadata,
	}
	updateResponse := objects.Update(command.Ctx.ServiceClient, containerName, objectName, updateOpts)
	if updateResponse.Err != nil {
		resource.Err = updateResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted metadata with keys [%s] from object [%s].\n", strings.Join(params.metadataKeys, ", "), objectName)
}
