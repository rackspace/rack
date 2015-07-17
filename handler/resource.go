package handler

import "github.com/jrperritt/rack/util"

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

// FlattenMap is used to flatten out a `map[string]map[string]interface{}`
func (resource *Resource) FlattenMap(key string) {
	keys := resource.Keys
	res := resource.Result.(map[string]interface{})
	if m, ok := res[key]; ok && util.Contains(keys, key) {
		for k, v := range m.(map[string]interface{}) {
			res[k] = v
			keys = append(keys, k)
		}
	}
	delete(res, key)
	resource.Keys = util.RemoveFromList(keys, key)
}
