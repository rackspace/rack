package securitygrouprulecommands

import (
	"github.com/jrperritt/rack/internal/github.com/fatih/structs"
	osSecurityGroupRules "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func securityGroupRuleSingle(rule *osSecurityGroupRules.SecGroupRule) map[string]interface{} {
	m := structs.Map(rule)

	m["SecurityGroupID"] = m["SecGroupID"]

	return m
}
