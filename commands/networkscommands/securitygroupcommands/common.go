package securitygroupcommands

import (
	"github.com/fatih/structs"
	osSecurityGroups "github.com/rackspace/gophercloud/openstack/networking/v2/extensions/security/groups"
	osSecurityGroupRules "github.com/rackspace/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func securityGroupSingle(securityGroup *osSecurityGroups.SecGroup) map[string]interface{} {
	m := structs.Map(securityGroup)

	tmpMap := make([]map[string]interface{}, len(m["Rules"].([]osSecurityGroupRules.SecGroupRule)))
	for i, pool := range m["Rules"].([]osSecurityGroupRules.SecGroupRule) {
		tmpMap[i] = map[string]interface{}{
			"ID":        pool.ID,
			"Direction": pool.Direction,
			"EtherType": pool.EtherType,
			"Protocol":  pool.Protocol,
		}
	}

	m["Rules"] = tmpMap

	/*
		if rules, ok := m["Rules"].([]osSecurityGroupRules.SecGroupRule); ok && len(rules) > 0 {
			out := []string{"ID\tDirection\tEtherType\tProtocol"}
			for _, rule := range rules {
				out = append(out, fmt.Sprintf("%s\t%s\t%s\t%s", rule.ID, rule.Direction, rule.EtherType, rule.Protocol))
			}
			m["Rules"] = strings.Join(out, "\n")
		} else {
			m["Rules"] = ""
		}
	*/

	return m
}
