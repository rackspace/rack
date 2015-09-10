package stacks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	"github.com/rackspace/rack/internal/gopkg.in/yaml.v2"
)

type TE struct {
	// Bin stores the contents of the template or environment.
	Bin []byte
	// URL stores the URL of the template. This is allowed to be a 'file://'
	// for local files.
	URL string
	// Parsed contains a parsed version of Bin. Since there are 2 different
	// fields referring to the same value, you must be careful when accessing
	// this filed.
	Parsed map[string]interface{}
	// Files contains a mapping between the urls in templates to their contents.
	Files map[string]string
	// fileMaps is a map used internally when determining Files.
	fileMaps map[string]string
	// baseURL represents the location of the template or environment file.
	baseURL string
}

func (t *TE) Fetch() error {
	if t.baseURL == "" {
		basePath, err := filepath.Abs(".")
		if err != nil {
			return err
		}
		u, err := gophercloud.NormalizePathURL("", basePath)
		if err != nil {
			return err
		}
		t.baseURL = u
	}
	if t.Bin != nil {
		// already have contents
		return nil
	}
	u, err := gophercloud.NormalizePathURL(t.baseURL, t.URL)
	if err != nil {
		return err
	}
	t.URL = u
	body, err := fetch(t.URL)
	if err != nil {
		return err
	}
	t.Bin = body
	return nil
}

// fetch the contents from URL, if not in Bin
func fetch(url string) ([]byte, error) {

	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c := &http.Client{Transport: transport}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// parse the contents and validate
func (t *TE) Parse() error {
	if err := t.Fetch(); err != nil {
		return err
	}
	if jerr := json.Unmarshal(t.Bin, &t.Parsed); jerr != nil {
		if yerr := yaml.Unmarshal(t.Bin, &t.Parsed); yerr != nil {
			return errors.New(fmt.Sprintf("Data in neither json nor yaml format."))
		}
	}
	return t.Validate()
}

// base Validate method, always returns true
func (t *TE) Validate() error {
	return nil
}

type igFunc func(string, interface{}) bool

func (t *TE) GetFileContents(te interface{}, ignoreIf igFunc) error {
	if t.Files == nil {
		t.Files = make(map[string]string)
	}
	if t.fileMaps == nil {
		t.fileMaps = make(map[string]string)
	}
	switch te.(type) {
	case map[string]interface{}, map[interface{}]interface{}:
		te_map, err := toStringKeys(te)
		if err != nil {
			return err
		}

		for k, v := range te_map {
			value, ok := v.(string)
			if !ok {
				if err := t.GetFileContents(v, ignoreIf); err != nil {
					return err
				}
			} else if !ignoreIf(k, value) {
				// at this point, the k, v pair has a reference to an external template.
				// The assumption of heatclient is that value v is a relative reference
				// to a file in the users environment
				childTE := new(TE)
				baseURL, err := gophercloud.NormalizePathURL(t.baseURL, value)
				if err != nil {
					return err
				}
				childTE.baseURL = baseURL
				if err := childTE.Parse(); err != nil {
					//TODO: log this error!
					return nil
				}
				if err := childTE.GetFileContents(childTE.Parsed, ignoreIf); err != nil {
					return err
				}
				// update parent template with current child templates' content
				t.fileMaps[value] = childTE.URL
				t.Files[childTE.URL] = string(childTE.Bin)
			}
		}
		return nil
	case []interface{}:
		te_slice := te.([]interface{})
		for i := range te_slice {
			if err := t.GetFileContents(te_slice[i], ignoreIf); err != nil {
				return err
			}
		}
	case string, bool, float64, nil, int:
		return nil
	default:
		return errors.New(fmt.Sprintf("%v: Unrecognized type", reflect.TypeOf(te)))

	}
	return nil
}

// convert map[interface{}]interface{} to map[string]interface{}
func toStringKeys(m interface{}) (map[string]interface{}, error) {
	switch m.(type) {
	case map[string]interface{}, map[interface{}]interface{}:
		typed_map := make(map[string]interface{})
		if _, ok := m.(map[interface{}]interface{}); ok {
			for k, v := range m.(map[interface{}]interface{}) {
				typed_map[k.(string)] = v
			}
		} else {
			typed_map = m.(map[string]interface{})
		}
		return typed_map, nil
	default:
		return nil, errors.New(fmt.Sprintf("Expected a map of type map[string]interface{} or map[interface{}]interface{}, actual type: %v", reflect.TypeOf(m)))

	}
}

// fix the template reference to files
func (t *TE) FixFileRefs() {
	t_str := string(t.Bin)
	for k, v := range t.fileMaps {
		t_str = strings.Replace(t_str, k, v, -1)
	}
	t.Bin = []byte(t_str)
}
