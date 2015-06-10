package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/rackspace/gophercloud"
)

// Creds are values that are cached to prevent re-authenticating for
// every CLI command.
type Creds struct {
	Ao     gophercloud.AuthOptions
	Region string
}

// GetCacheValue returns the cached value for the given key if it exists. It
func GetCacheValue(cacheKey string) (map[string]interface{}, error) {
	dir, err := rackDir()
	if err != nil {
		return nil, fmt.Errorf("Error reading from cache: %s", err)
	}
	f := path.Join(dir, "cache")
	cacheRaw, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("Error reading from cache: %s", err)
	}
	var m map[string]map[string]interface{}
	err = json.Unmarshal(cacheRaw, m)
	if err != nil {
		return nil, fmt.Errorf("Error reading from cache: %s", err)
	}
	return m[cacheKey], nil
}

// GetCacheValueCreds returns the user's cached credentials.
func GetCacheValueCreds(cacheKey string) (*Creds, error) {
	m, err := GetCacheValue(cacheKey)
	if err != nil {
		return nil, err
	}
	if m != nil {
		if raw := m["creds"]; raw != nil {
			if cc, ok := raw.(Creds); ok {
				return &cc, nil
			}
		}
	}
	return nil, nil
}

// SetCacheValueCreds writes the user's current credentials to the cache.
func SetCacheValueCreds(cacheKey, cacheValue *Creds) {

}

// GetCacheValueProviderClient returns the user's cached provider client.
func GetCacheValueProviderClient(cacheKey string) (*gophercloud.ProviderClient, error) {
	m, err := GetCacheValue(cacheKey)
	if err != nil {
		return nil, err
	}
	if m != nil {
		if raw := m["pc"]; raw != nil {
			if cc, ok := raw.(gophercloud.ProviderClient); ok {
				return &cc, nil
			}
		}
	}
	return nil, nil
}

// SetCacheValueProviderClient writes the user's current provider client to the cache.
func SetCacheValueProviderClient(cacheKey, cacheValue *gophercloud.ProviderClient) {

}
