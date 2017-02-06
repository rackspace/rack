package instancecommands

import (
	"fmt"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var deleteMetadata = cli.Command{
	Name:        "delete-metadata",
	Usage:       util.Usage(commandPrefix, "delete-metadata", ""),
	Description: "Delete metadata associated with the given server",
	Action:      actionDeleteMetadata,
	Flags:       commandoptions.CommandFlags(flagsDeleteMetadata, keysDeleteMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDeleteMetadata, keysDeleteMetadata))
	},
}

func flagsDeleteMetadata() []cli.Flag {
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
			Name:  "metadata-keys",
			Usage: "[required] A comma-separated string of keys of the metadata to delete from the server.",
		},
	}
}

var keysDeleteMetadata = []string{}

type paramsDeleteMetadata struct {
	serverID     string
	metadataKeys []string
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
	err := command.Ctx.CheckFlagsSet([]string{"metadata-keys"})
	if err != nil {
		return err
	}

	serverID, err := command.Ctx.IDOrName(osServers.IDFromName)
	resource.Params = &paramsDeleteMetadata{
		serverID:     serverID,
		metadataKeys: strings.Split(command.Ctx.CLIContext.String("metadata-keys"), ","),
	}
	return err
}

func (command *commandDeleteMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDeleteMetadata)
	for _, key := range params.metadataKeys {
		err := osServers.DeleteMetadatum(command.Ctx.ServiceClient, params.serverID, key).ExtractErr()
		if err != nil {
			resource.Err = err
			return
		}
	}
	resource.Result = fmt.Sprintf("Successfully deleted metadata")
}
