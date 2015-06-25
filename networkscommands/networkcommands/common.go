package networkcommands

import (
	"github.com/fatih/structs"
	osNetworks "github.com/rackspace/gophercloud/openstack/networking/v2/networks"
)

func networkSingle(network *osNetworks.Network) map[string]interface{} {
	m := structs.Map(network)
	m["Up"] = m["AdminStateUp"]
	m["Tenant ID"] = m["TenantID"]
	return m
}
