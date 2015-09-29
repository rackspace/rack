package instancecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osServers "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var updateMetadata = cli.Command{
	Name:        "update-metadata",
	Usage:       util.Usage(commandPrefix, "update-metadata", ""),
	Description: "Update metadata on the given server.",
	Action:      actionUpdateMetadata,
	Flags:       commandoptions.CommandFlags(flagsUpdateMetadata, keysUpdateMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsUpdateMetadata, keysUpdateMetadata))
	},
}

func flagsUpdateMetadata() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't provided] The server ID with the metadata.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `name` isn't provided] The server name with the metadata.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[required] A comma-separated string of 'key=value' pairs to update as metadata for the server.",
		},
	}
}

var keysUpdateMetadata = []string{}

type paramsUpdateMetadata struct {
	serverID string
	opts     *osServers.MetadataOpts
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
	metadata, err := command.Ctx.CheckKVFlag("metadata")
	if err != nil {
		return err
	}
	opts := osServers.MetadataOpts(metadata)

	serverID, err := command.Ctx.IDOrName(osServers.IDFromName)
	resource.Params = &paramsUpdateMetadata{
		serverID: serverID,
		opts:     &opts,
	}
	return err
}

func (command *commandUpdateMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsUpdateMetadata)
	metadata, err := osServers.UpdateMetadata(command.Ctx.ServiceClient, params.serverID, params.opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = metadata
}

func (command *commandUpdateMetadata) PreCSV(resource *handler.Resource) error {
	resource.Result = map[string]interface{}{
		"Metadata": resource.Result,
	}
	resource.Keys = []string{"Metadata"}
	resource.FlattenMap("Metadata")
	return nil
}

func (command *commandUpdateMetadata) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
