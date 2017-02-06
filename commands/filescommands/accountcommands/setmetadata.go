package accountcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osAccounts "github.com/rackspace/gophercloud/openstack/objectstorage/v1/accounts"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/accounts"
	"github.com/rackspace/rack/util"
)

var setMetadata = cli.Command{
	Name:        "set-metadata",
	Usage:       util.Usage(commandPrefix, "set-metadata", "--name <containerName> --metadata <metadata>"),
	Description: "Set metadata for the account. This will erase any current metadata.",
	Action:      actionSetMetadata,
	Flags:       commandoptions.CommandFlags(flagsSetMetadata, keysSetMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsSetMetadata, keysSetMetadata))
	},
}

func flagsSetMetadata() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[required] A comma-separated string of 'key=value' pairs to set as metadata for the account.",
		},
	}
}

var keysSetMetadata = []string{}

type paramsSetMetadata struct {
	metadata map[string]string
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
	err := command.Ctx.CheckFlagsSet([]string{"metadata"})
	if err != nil {
		return err
	}

	metadata, err := command.Ctx.CheckKVFlag("metadata")
	if err != nil {
		return err
	}

	resource.Params = &paramsSetMetadata{
		metadata: metadata,
	}
	return err
}

func (command *commandSetMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsSetMetadata)

	currentMetadata, err := accounts.Get(command.Ctx.ServiceClient).ExtractMetadata()
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

	updateOpts := osAccounts.UpdateOpts{
		Metadata:       params.metadata,
		DeleteMetadata: keys,
	}
	updateResponse := accounts.Update(command.Ctx.ServiceClient, updateOpts)
	if updateResponse.Err != nil {
		resource.Err = updateResponse.Err
		return
	}

	metadata, err := accounts.Get(command.Ctx.ServiceClient).ExtractMetadata()
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
