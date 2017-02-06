package stackeventcommands

import (
	"github.com/fatih/structs"
	osStackEvents "github.com/rackspace/gophercloud/openstack/orchestration/v1/stackevents"
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
		if event.Links != nil {
			links := make([]map[string]interface{}, len(event.Links))
			for i, link := range event.Links {
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
