package instancecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var listAddresses = cli.Command{
	Name:        "list-addresses",
	Usage:       util.Usage(commandPrefix, "list-addresses", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description: "Lists existing IP addresses for the server",
	Action:      actionListAddresses,
	Flags:       commandoptions.CommandFlags(flagsListAddresses, keysListAddresses),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsListAddresses, keysListAddresses))
	},
}

func flagsListAddresses() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The server ID from which to list the IP addresses.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The server name from which to list the IP addresses.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
	}
}

var keysListAddresses = []string{"Type", "Version", "Address"}

type paramsListAddresses struct {
	serverID string
}

type commandListAddresses handler.Command

func actionListAddresses(c *cli.Context) {
	command := &commandListAddresses{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandListAddresses) Context() *handler.Context {
	return command.Ctx
}

func (command *commandListAddresses) Keys() []string {
	return keysListAddresses
}

func (command *commandListAddresses) ServiceClientType() string {
	return serviceClientType
}

func (command *commandListAddresses) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsListAddresses{}
	return nil
}

func (command *commandListAddresses) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsListAddresses).serverID = item
	return nil
}

func (command *commandListAddresses) HandleSingle(resource *handler.Resource) error {
	id, err := command.Context().IDOrName(osServers.IDFromName)
	resource.Params.(*paramsListAddresses).serverID = id
	return err
}

func (command *commandListAddresses) Execute(resource *handler.Resource) {
	serverID := resource.Params.(*paramsListAddresses).serverID
	pager := servers.ListAddresses(command.Ctx.ServiceClient, serverID)
	var result []map[string]interface{}
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		addressesMap, err := servers.ExtractAddresses(page)
		if err != nil {
			return false, err
		}
		for t, addresses := range addressesMap {
			for _, address := range addresses {
				m := map[string]interface{}{
					"Type":    t,
					"Version": address.Version,
					"Address": address.Address,
				}
				result = append(result, m)
			}
		}
		return true, nil
	})
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = result
}

func (command *commandListAddresses) StdinField() string {
	return "id"
}
