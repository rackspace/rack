package accountcommands

import (
	"fmt"
	"strings"

	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	osAccounts "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/accounts"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/accounts"
	"github.com/jrperritt/rack/util"
)

var deleteMetadata = cli.Command{
	Name:        "delete-metadata",
	Usage:       util.Usage(commandPrefix, "delete-metadata", "--metadata-keys <metadataKeys>"),
	Description: "Delete specific metadata associated with the account.",
	Action:      actionDeleteMetadata,
	Flags:       commandoptions.CommandFlags(flagsDeleteMetadata, keysDeleteMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDeleteMetadata, keysDeleteMetadata))
	},
}

func flagsDeleteMetadata() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "metadata-keys",
			Usage: "[required] A comma-separated string of metadata keys to delete from the account.",
		},
	}
}

var keysDeleteMetadata = []string{}

type paramsDeleteMetadata struct {
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
	metadataKeys := strings.Split(command.Ctx.CLIContext.String("metadata-keys"), ",")
	for i, k := range metadataKeys {
		metadataKeys[i] = strings.Title(k)
	}

	resource.Params = &paramsDeleteMetadata{
		metadataKeys: metadataKeys,
	}
	return nil
}

func (command *commandDeleteMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDeleteMetadata)

	updateOpts := osAccounts.UpdateOpts{
		DeleteMetadata: params.metadataKeys,
	}
	updateResponse := accounts.Update(command.Ctx.ServiceClient, updateOpts)
	if updateResponse.Err != nil {
		resource.Err = updateResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted metadata with keys [%s] from account.\n", strings.Join(params.metadataKeys, ", "))
}
