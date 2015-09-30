package stackeventcommands

import (
	"encoding/json"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osStackEvents "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stackevents"
)

func eventSingle(rawEvent interface{}) map[string]interface{} {
	m := structs.Map(rawEvent)
	switch event := rawEvent.(type) {
	case *osStackEvents.Event:
		if event.Time.Unix() != -62135596800 {
			m["Time"] = event.Time
		} else {
			m["Time"] = ""
		}
		if resourceProperties, err := json.MarshalIndent(event.ResourceProperties, "", "  "); err != nil {
			m["ResourceProperties"] = ""
		} else {
			m["ResourceProperties"] = string(resourceProperties)
		}
		if event.Links != nil {
			links := make([]string, len(event.Links))
			for i, link := range event.Links {
				links[i] = link.PrettyPrintJSON()
			}
			m["Links"] = links
		}
		return m
	}
	return nil
}
