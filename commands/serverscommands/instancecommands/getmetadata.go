package instancecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var getMetadata = cli.Command{
	Name:        "get-metadata",
	Usage:       util.Usage(commandPrefix, "get-metadata", ""),
	Description: "Get all metadata associated with the given server",
	Action:      actionGetMetadata,
	Flags:       commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata))
	},
}

func flagsGetMetadata() []cli.Flag {
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
			Name:  "fields",
			Usage: "[optional] A comma-separated list of metadata keys to return.",
		},
	}
}

var keysGetMetadata = []string{}

type paramsGetMetadata struct {
	serverID string
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
	serverID, err := command.Ctx.IDOrName(osServers.IDFromName)
	resource.Params = &paramsGetMetadata{
		serverID: serverID,
	}
	return err
}

func (command *commandGetMetadata) Execute(resource *handler.Resource) {
	serverID := resource.Params.(*paramsGetMetadata).serverID
	metadata, err := servers.Metadata(command.Ctx.ServiceClient, serverID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = metadata
}

func (command *commandGetMetadata) PreCSV(resource *handler.Resource) error {
	resource.Result = map[string]interface{}{
		"Metadata": resource.Result,
	}
	resource.Keys = []string{"Metadata"}
	resource.FlattenMap("Metadata")
	return nil
}

func (command *commandGetMetadata) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
