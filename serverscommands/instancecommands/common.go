package instancecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/util"
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

func serverSingle(rawServer interface{}) map[string]interface{} {
	server, ok := rawServer.(*osServers.Server)
	if !ok {
		return nil
	}
	m := structs.Map(rawServer)
	m["Public IPv4"] = server.AccessIPv4
	m["Public IPv6"] = server.AccessIPv6
	m["Private IPv4"] = ""
	ips, ok := server.Addresses["private"].([]interface{})
	if ok || len(ips) > 0 {
		priv, ok := ips[0].(map[string]interface{})
		if ok {
			m["Private IPv4"] = priv["addr"]
		}
	}
	m["Flavor"] = server.Flavor["id"]
	m["Image"] = server.Image["id"]
	return m
}
