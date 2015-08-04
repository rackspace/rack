package accountcommands

import (
	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/accounts"
	"github.com/jrperritt/rack/util"
)

var getMetadata = cli.Command{
	Name:        "get-metadata",
	Usage:       util.Usage(commandPrefix, "get-metadata", ""),
	Description: "Get metadata associated with the account.",
	Action:      actionGetMetadata,
	Flags:       commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata))
	},
}

func flagsGetMetadata() []cli.Flag {
	return []cli.Flag{}
}

var keysGetMetadata = []string{"Metadata"}

type paramsGetMetadata struct {
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
	resource.Params = &paramsGetMetadata{}
	return nil
}

func (command *commandGetMetadata) Execute(resource *handler.Resource) {
	metadata, err := accounts.Get(command.Ctx.ServiceClient).ExtractMetadata()
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
