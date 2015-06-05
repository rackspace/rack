package servercommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rackcli/util"
	"github.com/rackspace/gophercloud"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func idOrName(c *cli.Context, client *gophercloud.ServiceClient) string {
	var err error
	var serverID string
	if c.IsSet("id") {
		serverID = c.String("id")
	} else if c.IsSet("name") {
		serverName := c.String("name")
		serverID, err = osServers.IDFromName(client, serverName)
		if err != nil {
			fmt.Printf("Error converting server name [%s] to ID: %s\n", serverName, err)
			os.Exit(1)
		}
	} else {
		util.PrintError(c, util.ErrMissingFlag{
			Msg: "One of either --id or --name must be provided.",
		})
	}
	return serverID
}

var idAndNameFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional; required if 'name' is not provided] The ID of the server to update",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if 'id' is not provided] The name of the server to update",
	},
}

var idOrNameUsage = "[--id <serverID> | --name <serverName>]"
