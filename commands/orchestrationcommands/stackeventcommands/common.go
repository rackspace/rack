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
			m["Time"] = "None"
		}
		if resourceProperties, err := json.MarshalIndent(event.ResourceProperties, "", "  "); err != nil {
			m["ResourceProperties"] = ""
		} else {
			m["ResourceProperties"] = string(resourceProperties)
		}
		return m
	}
	return nil
}
