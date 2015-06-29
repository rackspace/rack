package portcommands

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	osPorts "github.com/rackspace/gophercloud/openstack/networking/v2/ports"
)

func portSingle(port *osPorts.Port) map[string]interface{} {
	m := structs.Map(port)

	if fixedIPs, ok := m["FixedIPs"].([]osPorts.IP); ok && len(fixedIPs) > 0 {
		out := []string{"Subnet ID\tIP Address"}
		for _, ip := range fixedIPs {
			out = append(out, fmt.Sprintf("%s\t%s", ip.SubnetID, ip.IPAddress))
		}
		m["FixedIPs"] = strings.Join(out, "\n")
	} else {
		m["FixedIPs"] = ""
	}

	if nameServers, ok := m["SecurityGroups"].([]string); ok && len(nameServers) > 0 {
		m["SecurityGroups"] = strings.Join(nameServers, "\n")
	} else {
		m["SecurityGroups"] = ""
	}

	return m
}
