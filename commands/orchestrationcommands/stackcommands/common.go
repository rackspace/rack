package stackcommands

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/rackspace/gophercloud"
	osStacks "github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
)

func stackList(client *gophercloud.ServiceClient) ([]map[string]interface{}, error) {
	pager := osStacks.List(client, nil)
	pages, err := pager.AllPages()
	if err != nil {
		return nil, err
	}
	info, err := osStacks.ExtractStacks(pages)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, len(info))
	for j, stack := range info {
		result[j] = structs.Map(&stack)
		result[j]["CreationTime"] = stack.CreationTime
		result[j]["UpdatedTime"] = stack.UpdatedTime
	}
	return result, nil
}

func stackSingle(rawStack interface{}) map[string]interface{} {
	m := structs.Map(rawStack)
	switch stack := rawStack.(type) {
	case *osStacks.ListedStack:
		m["CreationTime"] = stack.CreationTime
		if stack.UpdatedTime.Unix() != -62135596800 {
			m["UpdatedTime"] = stack.UpdatedTime
		} else {
			m["UpdatedTime"] = ""
		}
		if stack.Links != nil {
			links := make([]map[string]interface{}, len(stack.Links))
			for i, link := range stack.Links {
				links[i] = map[string]interface{}{
					"Href": link.Href,
					"Rel":  link.Rel,
				}
			}
			m["Links"] = links
		}
		return m
	case *osStacks.CreatedStack:
		if stack.Links != nil {
			links := make([]map[string]interface{}, len(stack.Links))
			for i, link := range stack.Links {
				links[i] = map[string]interface{}{
					"Href": link.Href,
					"Rel":  link.Rel,
				}
			}
			m["Links"] = links
		}
		return m
	case *osStacks.RetrievedStack:
		m["CreationTime"] = stack.CreationTime
		if stack.UpdatedTime.Unix() != -62135596800 {
			m["UpdatedTime"] = stack.UpdatedTime
		} else {
			m["UpdatedTime"] = ""
		}
		if stack.Links != nil {
			links := make([]map[string]interface{}, len(stack.Links))
			for i, link := range stack.Links {
				links[i] = map[string]interface{}{
					"Href": link.Href,
					"Rel":  link.Rel,
				}
			}
			m["Links"] = links
		}
		return m
	case *osStacks.PreviewedStack:
		m["CreationTime"] = stack.CreationTime
		if stack.UpdatedTime.Unix() != -62135596800 {
			m["UpdatedTime"] = stack.UpdatedTime
		} else {
			m["UpdatedTime"] = ""
		}
		if stack.Links != nil {
			links := make([]map[string]interface{}, len(stack.Links))
			for i, link := range stack.Links {
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

// IDAndName is a function that returns both id and name when either one is specified
func IDAndName(client *gophercloud.ServiceClient, name, id string) (string, string, error) {
	// no checking is done when both name and id are provided
	if name != "" && id != "" {
		return name, id, nil
	} else if name == "" && id == "" {
		return "", "", fmt.Errorf("Either name or id must be provided.")
	}
	result, err := stackList(client)
	if err != nil {
		return "", "", err
	}
	if name == "" {
		for _, stack := range result {
			if stack["ID"] == id {
				return stack["Name"].(string), id, nil
			}
		}
	} else {
		for _, stack := range result {
			if stack["Name"] == name {
				return name, stack["ID"].(string), nil
			}
		}
	}
	return name, id, fmt.Errorf("Stack with name: %s and id: %s not found.", name, id)
}
