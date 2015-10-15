package stackresourcecommands

import (
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osStackResources "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stackresources"
)

func resourceSingle(rawResource interface{}) map[string]interface{} {
	m := structs.Map(rawResource)
	switch resource := rawResource.(type) {
	case *osStackResources.Resource:
		if resource.CreationTime.Unix() != -62135596800 {
			m["CreationTime"] = resource.CreationTime
		} else {
			m["CreationTime"] = ""
		}
		if resource.UpdatedTime.Unix() != -62135596800 {
			m["UpdatedTime"] = resource.UpdatedTime
		} else {
			m["UpdatedTime"] = ""
		}
		if resource.Links != nil {
			links := make([]map[string]interface{}, len(resource.Links))
			for i, link := range resource.Links {
				links[i] = map[string]interface{}{
					"Href": link.Href,
					"Rel":  link.Rel,
				}
			}
			m["Links"] = links
		}
		return m
	}
	return nil

}
