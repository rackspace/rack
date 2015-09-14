package stackcommands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
)

func stackList(client *gophercloud.ServiceClient) ([]map[string]interface{}, error) {
	pager := stacks.List(client, nil)
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
		// TODO: fix the decoding/parsing to make this work right
		result[j]["CreationTime"] = stack.CreationTime
		result[j]["UpdatedTime"] = stack.UpdatedTime
	}
	return result, nil
}

func stackSingle(rawStack interface{}) map[string]interface{} {
	m := structs.Map(rawStack)
	switch stack := rawStack.(type) {
	case *osStacks.RetrievedStack:
		m["CreationTime"] = stack.CreationTime
		if stack.UpdatedTime.Unix() != -62135596800 {
			m["UpdatedTime"] = stack.UpdatedTime
		} else {
			m["UpdatedTime"] = "None"
		}
		if parameters, err := json.MarshalIndent(stack.Parameters, "", "  "); err != nil {
			m["Parameters"] = ""
		} else {
			m["Parameters"] = string(parameters)
		}
		if outputs, err := json.MarshalIndent(stack.Outputs, "", "  "); err != nil {
			m["Outputs"] = ""
		} else {
			m["Outputs"] = string(outputs)
		}
		if stack.Timeout == 0 {
			m["Timeout"] = "None"
		}
		if stack.Links != nil {
			links := make([]string, len(stack.Links))
			for i, link := range stack.Links {
				links[i] = link.PrettyPrintJSON()
			}
			m["Links"] = links
		}
		return m
	case *osStacks.PreviewedStack:
		m["CreationTime"] = stack.CreationTime
		if stack.UpdatedTime.Unix() != -62135596800 {
			m["UpdatedTime"] = stack.UpdatedTime
		} else {
			m["UpdatedTime"] = "None"
		}
		if parameters, err := json.MarshalIndent(stack.Parameters, "", "  "); err != nil {
			m["Parameters"] = ""
		} else {
			m["Parameters"] = string(parameters)
		}
		if resourceJSON, err := json.MarshalIndent(stack.Resources, "", "  "); err != nil {
			m["Resources"] = "None"
		} else {
			m["Resources"] = string(resourceJSON)
		}
		if stack.Timeout == 0 {
			m["Timeout"] = "None"
		}
		if stack.Links != nil {
			links := make([]string, len(stack.Links))
			for i, link := range stack.Links {
				links[i] = link.PrettyPrintJSON()
			}
			m["Links"] = links
		}
		return m
	case *osStacks.AbandonedStack:
		if environment, err := json.MarshalIndent(stack.Environment, "", "  "); err != nil {
			m["Environment"] = ""
		} else {
			m["Environment"] = string(environment)
		}
		if files, err := json.MarshalIndent(stack.Files, "", "  "); err != nil {
			m["Files"] = ""
		} else {
			m["Files"] = string(files)
		}
		if resources, err := json.MarshalIndent(stack.Resources, "", "  "); err != nil {
			m["Resources"] = ""
		} else {
			m["Resources"] = string(resources)
		}
		if template, err := json.MarshalIndent(stack.Template, "", "  "); err != nil {
			m["Template"] = ""
		} else {
			m["Template"] = string(template)
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
		return "", "", errors.New("Either name or id must be provided.")
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
	return name, id, errors.New(fmt.Sprintf("Stack with name: %s and id: %s not found.", name, id))
}
