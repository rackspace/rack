package instancecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
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

func mungeServerMap(m map[string]interface{}) {
	m["Public IPv4"] = m["AccessIPv4"]
	m["Public IPv6"] = m["AccessIPv6"]

	// Private IPv4 case
	// Nested m["Addresses"]["private"][0]
	addrs, ok := m["Addresses"].(map[string]interface{})
	if ok {
		ips, ok := addrs["private"].([]interface{})
		if ok || len(ips) > 0 {
			priv, ok := ips[0].(map[string]interface{})
			if ok {
				m["Private IPv4"] = priv["addr"]
			}
		}
	}
	if !ok { // if any were not ok, set the field to blank
		m["Private IPv4"] = ""
	}

	flavor, ok := m["Flavor"].(map[string]interface{})
	if ok {
		m["Flavor"] = flavor["id"]
	} else {
		m["Flavor"] = ""
	}
	image, ok := m["Image"].(map[string]interface{})
	if ok {
		m["Image"] = image["id"]
	} else {
		m["Image"] = ""
	}
}
