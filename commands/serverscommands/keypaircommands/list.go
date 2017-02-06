package keypaircommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
	"github.com/rackspace/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists keypairs",
	Action:      actionList,
	Flags:       commandoptions.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{}
}

var keysList = []string{"Name", "Fingerprint"}

type paramsList struct{}

type commandList handler.Command

func actionList(c *cli.Context) {
	command := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandList) Context() *handler.Context {
	return command.Ctx
}

func (command *commandList) Keys() []string {
	return keysList
}

func (command *commandList) ServiceClientType() string {
	return serviceClientType
}

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	err := keypairs.List(command.Ctx.ServiceClient).EachPage(func(page pagination.Page) (bool, error) {
		info, err := osKeypairs.ExtractKeyPairs(page)
		if err != nil {
			return false, err
		}
		result := make([]map[string]interface{}, len(info))
		for j, key := range info {
			result[j] = structs.Map(key)
		}
		resource.Result = result
		return false, nil
	})
	if err != nil {
		resource.Err = err
		return
	}
}
