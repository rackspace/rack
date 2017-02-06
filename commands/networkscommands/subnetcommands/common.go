package subnetcommands

import (
	"strings"

	"github.com/fatih/structs"
	osSubnets "github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
)

func subnetSingle(subnet *osSubnets.Subnet) map[string]interface{} {
	m := structs.Map(subnet)

	tmpMap := make([]map[string]interface{}, len(m["AllocationPools"].([]osSubnets.AllocationPool)))
	for i, pool := range m["AllocationPools"].([]osSubnets.AllocationPool) {
		tmpMap[i] = structs.Map(pool)
	}

	m["AllocationPools"] = tmpMap

	/*
		m["Host Routes"] = m["HostRoutes"]
		delete(m, "HostRoutes")
		if hostRoutes, ok := m["Host Routes"].([]osSubnets.HostRoute); ok && len(hostRoutes) > 0 {
			fmt.Printf("hostRoutes: %+v\n", hostRoutes)
			out := []string{"Destination CIDR\tNext Hop"}
			for _, route := range hostRoutes {
				out = append(out, fmt.Sprintf("%s\t%s", route.DestinationCIDR, route.NextHop))
			}
			m["Host Routes"] = strings.Join(out, "\n")
		} else {
			m["Host Routes"] = ""
		}
	*/

	if nameServers, ok := m["DNSNameservers"].([]string); ok && len(nameServers) > 0 {
		m["DNSNameservers"] = strings.Join(nameServers, ", ")
	} else {
		m["DNSNameservers"] = ""
	}

	return m
}
