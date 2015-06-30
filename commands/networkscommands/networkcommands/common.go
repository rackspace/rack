package networkcommands

import (
	"strings"

	"github.com/jrperritt/rack/internal/github.com/fatih/structs"
	osNetworks "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/networking/v2/networks"
)

func networkSingle(network *osNetworks.Network) map[string]interface{} {
	m := structs.Map(network)
	m["Up"] = m["AdminStateUp"]
	m["Tenant ID"] = m["TenantID"]
	if subnets, ok := m["Subnets"].([]string); ok {
		m["Subnets"] = strings.Join(subnets, ",")
	}
	return m
}
