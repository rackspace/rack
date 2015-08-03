package containercommands

import (
	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	"github.com/jrperritt/rack/util"
)

var setMetadata = cli.Command{
	Name:        "set-metadata",
	Usage:       util.Usage(commandPrefix, "set-metadata", "--name <containerName> --metadata <metadata>"),
	Description: "Set metadata for the given container. This will erase any current metadata.",
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
			Usage: "[required] The container name with the metadata.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[required] A comma-separated string of 'key=value' pairs to set as metadata for the container.",
		},
	}
}

var keysSetMetadata = []string{}

type paramsSetMetadata struct {
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
	err := command.Ctx.CheckFlagsSet([]string{"name", "metadata"})
	if err != nil {
		return err
	}

	metadata, err := command.Ctx.CheckKVFlag("metadata")
	if err != nil {
		return err
	}

	resource.Params = &paramsSetMetadata{
		containerName: command.Ctx.CLIContext.String("name"),
		metadata:      metadata,
	}
	return err
}

func (command *commandSetMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsSetMetadata)
	containerName := params.containerName

	currentMetadata, err := containers.Get(command.Ctx.ServiceClient, containerName).ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}

	i := 0
	keys := make([]string, len(currentMetadata))
	for k := range currentMetadata {
		keys[i] = k
		i++
	}

	updateOpts := containers.UpdateOpts{
		Metadata:       params.metadata,
		DeleteMetadata: keys,
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
