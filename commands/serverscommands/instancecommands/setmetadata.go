package instancecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osServers "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var setMetadata = cli.Command{
	Name:        "set-metadata",
	Usage:       util.Usage(commandPrefix, "set-metadata", ""),
	Description: "Set metadata for the given server. This will erase any current metadata.",
	Action:      actionSetMetadata,
	Flags:       commandoptions.CommandFlags(flagsSetMetadata, keysSetMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsSetMetadata, keysSetMetadata))
	},
}

func flagsSetMetadata() []cli.Flag {
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
			Usage: "[required] A comma-separated string of 'key=value' pairs to set as metadata for the server.",
		},
	}
}

var keysSetMetadata = []string{}

type paramsSetMetadata struct {
	serverID string
	opts     *osServers.MetadataOpts
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
	metadata, err := command.Ctx.CheckKVFlag("metadata")
	if err != nil {
		return err
	}
	opts := osServers.MetadataOpts(metadata)

	serverID, err := command.Ctx.IDOrName(osServers.IDFromName)
	resource.Params = &paramsSetMetadata{
		serverID: serverID,
		opts:     &opts,
	}
	return err
}

func (command *commandSetMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsSetMetadata)
	metadata, err := osServers.ResetMetadata(command.Ctx.ServiceClient, params.serverID, params.opts).Extract()
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
