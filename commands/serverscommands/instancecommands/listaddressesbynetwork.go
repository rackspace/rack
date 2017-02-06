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

var listAddressesByNetwork = cli.Command{
	Name:        "list-addresses-by-network",
	Usage:       util.Usage(commandPrefix, "list-addresses-by-network", "--network <networkType> [--id <serverID> | --name <serverName> | --stdin id]"),
	Description: "Lists existing IP addresses for the given server and network",
	Action:      actionListAddressesByNetwork,
	Flags:       commandoptions.CommandFlags(flagsListAddressesByNetwork, keysListAddressesByNetwork),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsListAddressesByNetwork, keysListAddressesByNetwork))
	},
}

func flagsListAddressesByNetwork() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "network",
			Usage: "[required] The network for which to list the IP addresses.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The server ID from which to list the IP addresses.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The server name from which to list the IP addresses.",
		},
	}
}

var keysListAddressesByNetwork = []string{"Type", "Version", "Address"}

type paramsListAddressesByNetwork struct {
	serverID string
	network  string
}

type commandListAddressesByNetwork handler.Command

func actionListAddressesByNetwork(c *cli.Context) {
	command := &commandListAddressesByNetwork{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandListAddressesByNetwork) Context() *handler.Context {
	return command.Ctx
}

func (command *commandListAddressesByNetwork) Keys() []string {
	return keysListAddressesByNetwork
}

func (command *commandListAddressesByNetwork) ServiceClientType() string {
	return serviceClientType
}

func (command *commandListAddressesByNetwork) HandleFlags(resource *handler.Resource) error {
	serverID, err := command.Ctx.IDOrName(osServers.IDFromName)
	if err != nil {
		return err
	}

	err = command.Ctx.CheckFlagsSet([]string{"network"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext
	resource.Params = &paramsListAddressesByNetwork{
		serverID: serverID,
		network:  c.String("network"),
	}
	return nil
}

func (command *commandListAddressesByNetwork) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsListAddressesByNetwork)
	pager := servers.ListAddressesByNetwork(command.Ctx.ServiceClient, params.serverID, params.network)
	var result []map[string]interface{}
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		addresses, err := servers.ExtractNetworkAddresses(page)
		if err != nil {
			return false, err
		}
		for _, address := range addresses {
			m := map[string]interface{}{
				"Type":    params.network,
				"Version": address.Version,
				"Address": address.Address,
			}
			result = append(result, m)
		}
		return true, nil
	})
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = result
}
