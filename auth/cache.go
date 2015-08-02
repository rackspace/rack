package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud"
	"github.com/jrperritt/rack/util"
)

// Cache represents a place to store user authentication credentials.
type Cache struct {
	items map[string]CacheItem
	sync.RWMutex
}

// CacheItem represents a single item in the cache.
type CacheItem struct {
	TokenID         string
	ServiceEndpoint string
}

// CacheKey returns the cache key formed from the user's authentication credentials.
func CacheKey(ao gophercloud.AuthOptions, region, serviceClientType string, urlType gophercloud.Availability) string {
	if ao.Username != "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", ao.Username, ao.IdentityEndpoint, region, serviceClientType, urlType)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s", ao.TenantID, ao.IdentityEndpoint, region, serviceClientType, urlType)
}

func cacheFile() (string, error) {
	dir, err := util.RackDir()
	if err != nil {
		return "", fmt.Errorf("Error reading from cache: %s", err)
	}
	filepath := path.Join(dir, "cache")
	// check if the cache file exists
	if _, err := os.Stat(filepath); err == nil {
		return filepath, nil
	}
	// create the cache file if it doesn't already exist
	f, err := os.Create(filepath)
	defer f.Close()
	return filepath, err
}

// all returns all the values in the cache
func (cache *Cache) all() error {
	filename, err := cacheFile()
	if err != nil {
		return err
	}
	cache.RLock()
	defer cache.RUnlock()
	data, _ := ioutil.ReadFile(filename)
	if len(data) == 0 {
		cache.items = make(map[string]CacheItem)
		return nil
	}
	err = json.Unmarshal(data, &cache.items)
	if err != nil {
		return err
	}

	return nil
}

// Value returns the cached value for the given key if it exists.
func (cache *Cache) Value(cacheKey string) (*CacheItem, error) {
	err := cache.all()
	if err != nil {
		return nil, fmt.Errorf("Error getting cache value: %s", err)
	}
	creds := cache.items[cacheKey]
	if creds.TokenID == "" {
		return nil, nil
	}
	return &creds, nil
}

// SetValue writes the user's current provider client to the cache.
func (cache *Cache) SetValue(cacheKey string, cacheValue *CacheItem) error {
	// get cache items
	err := cache.all()
	if err != nil {
		return err
	}
	if cacheValue == nil {
		delete(cache.items, cacheKey)
	} else {
		// set cache value for cacheKey
		cache.items[cacheKey] = *cacheValue
	}
	filename, err := cacheFile()
	if err != nil {
		return err
	}
	cache.Lock()
	defer cache.Unlock()
	data, err := json.Marshal(cache.items)
	// write cache to file
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("Error setting cache value: %s", err)
	}
	return nil
}
