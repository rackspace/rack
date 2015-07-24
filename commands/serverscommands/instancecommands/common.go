package instancecommands

import (
	"github.com/jrperritt/rack/internal/github.com/fatih/structs"
	osServers "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func serverSingle(rawServer interface{}) map[string]interface{} {
	server, ok := rawServer.(*osServers.Server)
	if !ok {
		return nil
	}
	m := structs.Map(rawServer)
	m["PublicIPv4"] = server.AccessIPv4
	m["PublicIPv6"] = server.AccessIPv6
	m["PrivateIPv4"] = ""
	ips, ok := server.Addresses["private"].([]interface{})
	if ok || len(ips) > 0 {
		priv, ok := ips[0].(map[string]interface{})
		if ok {
			m["PrivateIPv4"] = priv["addr"]
		}
	}
	m["Flavor"] = server.Flavor["id"]
	m["Image"] = server.Image["id"]
	return m
}
