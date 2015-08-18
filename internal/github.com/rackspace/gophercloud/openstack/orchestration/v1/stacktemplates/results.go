package stacktemplates

import (
	"encoding/json"

	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
)

// GetResult represents the result of a Get operation.
type GetResult struct {
	gophercloud.Result
}

// Extract returns the JSON template and is called after a Get operation.
func (r GetResult) Extract() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	template, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return nil, err
	}
	return template, nil
}

// ValidateResult represents the result of a Validate operation.
type ValidateResult struct {
	gophercloud.Result
}

// Extract returns the JSON template and is called after a Get operation.
func (r ValidateResult) Extract() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	validateResult, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return nil, err
	}

	return validateResult, nil
}
