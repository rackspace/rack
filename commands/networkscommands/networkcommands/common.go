package networkcommands

import (
	"strings"

	"github.com/fatih/structs"
	osNetworks "github.com/rackspace/gophercloud/openstack/networking/v2/networks"
)

func networkSingle(network *osNetworks.Network) map[string]interface{} {
	m := structs.Map(network)
	m["Up"] = m["AdminStateUp"]
	if subnets, ok := m["Subnets"].([]string); ok {
		m["Subnets"] = strings.Join(subnets, ",")
	}
	return m
}
