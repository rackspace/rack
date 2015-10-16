package handler

import (
	"fmt"

	"github.com/rackspace/rack/util"
)

// Resource is a general resource from Rackspace. This object stores information
// about a single request and response from Rackspace.
type Resource struct {
	// Keys are the fields available to output. These may be limited by the `fields`
	// flag.
	Keys []string
	// Params will be the command-specific parameters, such as an instance ID or
	// list options.
	Params interface{}
	// Result will store the result of a single command.
	Result interface{}
	// Err will store any error encountered while processing the command.
	Err error
}

// FlattenMap is used to flatten out a `map[string]map[string]*`
func (resource *Resource) FlattenMap(key string) {

	res := resource.Result.(map[string]interface{})
	if m, ok := res[key]; ok && util.Contains(resource.Keys, key) {
		switch m.(type) {
		case []map[string]interface{}:
			for i, hashmap := range m.([]map[string]interface{}) {
				for k, v := range hashmap {
					newKey := fmt.Sprintf("%s%d:%s", key, i, k)
					res[newKey] = v
					resource.Keys = append(resource.Keys, newKey)
					resource.FlattenMap(newKey)
				}
			}
		case []interface{}:
			for i, element := range m.([]interface{}) {
				newKey := fmt.Sprintf("%s%d", key, i)
				res[newKey] = element
				resource.Keys = append(resource.Keys, newKey)
				resource.FlattenMap(newKey)
			}
		case map[string]interface{}, map[interface{}]interface{}:
			mMap := toStringKeys(m)
			for k, v := range mMap {
				newKey := fmt.Sprintf("%s:%s", key, k)
				res[newKey] = v
				resource.Keys = append(resource.Keys, newKey)
				resource.FlattenMap(newKey)
			}
		case map[string]string:
			for k, v := range m.(map[string]string) {
				newKey := fmt.Sprintf("%s:%s", key, k)
				res[newKey] = v
				resource.Keys = append(resource.Keys, newKey)
			}
		default:
			return
		}
		delete(res, key)
		resource.Keys = util.RemoveFromList(resource.Keys, key)
	}
}

// convert map[interface{}]interface{} to map[string]interface{}
func toStringKeys(m interface{}) map[string]interface{} {
	switch m.(type) {
	case map[interface{}]interface{}:
		typedMap := make(map[string]interface{})
		for k, v := range m.(map[interface{}]interface{}) {
			typedMap[k.(string)] = v
		}
		return typedMap
	case map[string]interface{}:
		typedMap := m.(map[string]interface{})
		return typedMap
	}
	return nil
}
