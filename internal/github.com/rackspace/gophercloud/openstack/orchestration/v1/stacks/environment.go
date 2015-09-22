package stacks

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Environment struct {
	TE
}

var EnvironmentSections = map[string]bool{
	"parameters":         true,
	"parameter_defaults": true,
	"resource_registry":  true,
}

func (e *Environment) Validate() error {
	if e.Parsed == nil {
		if err := e.Parse(); err != nil {
			return err
		}
	}
	for key, _ := range e.Parsed {
		if _, ok := EnvironmentSections[key]; !ok {
			return errors.New(fmt.Sprintf("Environment has wrong section: %s", key))
		}
	}
	return nil
}

// Parse environment file to resolve the urls of the resources
func GetRRFileContents(e *Environment, ignoreIf igFunc) error {
	if e.Files == nil {
		e.Files = make(map[string]string)
	}
	if e.fileMaps == nil {
		e.fileMaps = make(map[string]string)
	}
	rr := e.Parsed["resource_registry"]
	switch rr.(type) {
	case map[string]interface{}, map[interface{}]interface{}:
		rr_map, err := toStringKeys(rr)
		if err != nil {
			return err
		}
		var baseURL string
		if val, ok := rr_map["base_url"]; ok {
			baseURL = val.(string)
		} else {
			baseURL = e.baseURL
		}
		tempTemplate := new(TE)
		tempTemplate.baseURL = baseURL
		if err := tempTemplate.GetFileContents(rr, ignoreIf); err != nil {
			return err
		}
		if val, ok := rr_map["resources"]; ok {
			switch val.(type) {
			case map[string]interface{}, map[interface{}]interface{}:
				resources_map, err := toStringKeys(val)
				if err != nil {
					return err
				}
				for _, v := range resources_map {
					switch v.(type) {
					case map[string]interface{}, map[interface{}]interface{}:
						resource_map, err := toStringKeys(v)
						if err != nil {
							return err
						}
						var resourceBaseURL string
						if val, ok := resource_map["base_url"]; ok {
							resourceBaseURL = val.(string)
						} else {
							resourceBaseURL = baseURL
						}
						tempTemplate.baseURL = resourceBaseURL
						if err := tempTemplate.GetFileContents(v, ignoreIf); err != nil {
							return err
						}
					}

				}

			}
		}
		e.Files = tempTemplate.Files
		return nil
	default:
		return nil
	}
}

// function to choose keys whose values are other environment files
func ignoreIfEnvironment(key string, value interface{}) bool {
	if key == "base_url" || key == "hooks" {
		return true
	}
	if reflect.ValueOf(value).Kind() == reflect.Map {
		return true
	}
	valueString, ok := value.(string)
	if !ok {
		return true
	}
	if strings.Contains(valueString, "::") {
		return true
	}
	return false
}
