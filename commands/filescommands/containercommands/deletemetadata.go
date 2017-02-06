package containercommands

import (
	"fmt"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	"github.com/rackspace/rack/util"
)

var deleteMetadata = cli.Command{
	Name:        "delete-metadata",
	Usage:       util.Usage(commandPrefix, "delete-metadata", "--name <containerName> --metadata-keys <metadataKeys>"),
	Description: "Delete specific metadata from the given container.",
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
			Usage: "[required] The container name with the metadata.",
		},
		cli.StringFlag{
			Name:  "metadata-keys",
			Usage: "[required] A comma-separated string of metadata keys to delete from the container.",
		},
	}
}

var keysDeleteMetadata = []string{}

type paramsDeleteMetadata struct {
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
	err := command.Ctx.CheckFlagsSet([]string{"name", "metadata-keys"})
	if err != nil {
		return err
	}
	metadataKeys := strings.Split(command.Ctx.CLIContext.String("metadata-keys"), ",")
	for i, k := range metadataKeys {
		metadataKeys[i] = strings.Title(k)
	}

	resource.Params = &paramsDeleteMetadata{
		containerName: command.Ctx.CLIContext.String("name"),
		metadataKeys:  metadataKeys,
	}
	return nil
}

func (command *commandDeleteMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDeleteMetadata)
	containerName := params.containerName

	getResponse := containers.Get(command.Ctx.ServiceClient, containerName)
	if getResponse.Err != nil {
		resource.Err = getResponse.Err
		return
	}

	updateOpts := containers.UpdateOpts{
		DeleteMetadata: params.metadataKeys,
	}
	updateResponse := containers.Update(command.Ctx.ServiceClient, containerName, updateOpts)
	if updateResponse.Err != nil {
		resource.Err = updateResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted metadata with keys [%s] from container [%s].\n", strings.Join(params.metadataKeys, ", "), containerName)
}
