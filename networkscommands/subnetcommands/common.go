package subnetcommands

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	osSubnets "github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
)

func subnetSingle(subnet *osSubnets.Subnet) map[string]interface{} {
	m := structs.Map(subnet)

	if allocationPools, ok := m["AllocationPools"].([]osSubnets.AllocationPool); ok && len(allocationPools) > 0 {
		out := []string{"Start\tEnd"}
		for _, pool := range allocationPools {
			out = append(out, fmt.Sprintf("%s\t%s", pool.Start, pool.End))
		}
		m["Allocation Pools"] = strings.Join(out, "\n")
	}

	m["Host Routes"] = m["HostRoutes"]
	delete(m, "HostRoutes")
	if hostRoutes, ok := m["Host Routes"].([]osSubnets.HostRoute); ok && len(hostRoutes) > 0 {
		out := []string{"Destination CIDR\tNext Hop"}
		for _, route := range hostRoutes {
			out = append(out, fmt.Sprintf("%s\t%s", route.DestinationCIDR, route.NextHop))
		}
		m["Host Routes"] = strings.Join(out, "\n")
	} else {
		m["Host Routes"] = ""
	}

	if nameServers, ok := m["DNSNameservers"].([]string); ok && len(nameServers) > 0 {
		m["DNS Nameservers"] = strings.Join(nameServers, "\n")
	} else {
		m["DNS Nameservers"] = ""
	}

	m["Tenant ID"] = m["TenantID"]
	m["Network ID"] = m["NetworkID"]
	m["Gateway IP"] = m["GatewayIP"]

	return m
}
