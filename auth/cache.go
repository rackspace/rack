package auth

import (
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud"
)

// Cache represents a place to store user authentication credentials.
type Cache struct {
	items map[string]CacheItem
	os.File
	sync.RWMutex
}

// CacheItem represents a single item in the cache.
type CacheItem struct {
	IdentityBase     string
	IdentityEndpoint string
	TokenID          string
	HTTPClient       http.Client
	ServiceEndpoint  string
}

// CacheKey returns the cache key formed from the user's authentication credentials.
func CacheKey(ao gophercloud.AuthOptions, region string, serviceClientType string) string {
	return fmt.Sprintf("%s,%s,%s,%s", ao.Username, ao.IdentityEndpoint, region, serviceClientType)
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
	f, _ := os.Open(filename)
	defer f.Close()

	// gob is used instead of JSON because Go can't handle marshalling an
	// http.Client to a JSON object.
	err = gob.NewDecoder(f).Decode(&cache.items)
	if err != nil {
		if err == io.EOF {
			cache.items = make(map[string]CacheItem)
			return nil
		}
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
	// set cache value for cacheKey
	cache.items[cacheKey] = *cacheValue
	filename, err := cacheFile()
	if err != nil {
		return err
	}
	cache.Lock()
	defer cache.Unlock()
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error setting cache value: %s", err)
	}
	defer f.Close()
	// write cache to file
	err = gob.NewEncoder(f).Encode(cache.items)
	if err != nil {
		return fmt.Errorf("Error setting cache value: %s", err)
	}
	return nil
}
