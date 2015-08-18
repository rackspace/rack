package stackresourcecommands

import(
    "encoding/json"
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
            m["CreationTime"] = "None"
        }
        if resource.UpdatedTime.Unix() != -62135596800 {
            m["UpdatedTime"] = resource.UpdatedTime
        } else {
            m["UpdatedTime"] = "None"
        }
        if attributes, err := json.MarshalIndent(resource.Attributes, "", "  "); err != nil {
            m["Attributes"] = ""
        } else {
            m["Attributes"] = string(attributes)
        }
        if resource.Links != nil {
            links := make([]string, len(resource.Links))
            for i, link := range(resource.Links) {
                links[i] = link.PrettyPrintJSON()
            }
            m["Links"] = links
        }
        return m
    case *osStackResources.TypeSchema:
        if attributes, err := json.MarshalIndent(resource.Attributes, "", "  "); err != nil {
            m["Attributes"] = ""
        } else {
            m["Attributes"] = string(attributes)
        }
        if properties, err := json.MarshalIndent(resource.Properties, "", "  "); err != nil {
            m["Properties"] = ""
        } else {
            m["Properties"] = string(properties)
        }
        if support, err := json.MarshalIndent(resource.SupportStatus, "", "  "); err != nil {
            m["SupportStatus"] = ""
        } else {
            m["SupportStatus"] = string(support)
        }
        return m
    case *osStackResources.TypeTemplate:
        if outputs, err := json.MarshalIndent(resource.Outputs, "", "  "); err != nil {
            m["Outputs"] = ""
        } else {
            m["Outputs"] = string(outputs)
        }
        if parameters, err := json.MarshalIndent(resource.Parameters, "", "  "); err != nil {
            m["Parameters"] = ""
        } else {
            m["Parameters"] = string(parameters)
        }
        if resources, err := json.MarshalIndent(resource.Resources, "", "  "); err != nil {
            m["Resources"] = ""
        } else {
            m["Resources"] = string(resources)
        }
        return m
    }
    return nil

}
