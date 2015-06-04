package servercommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
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
		fmt.Printf("Either the 'id' or 'name' flag must be provided.\n")
		os.Exit(1)
	}
	return serverID
}
