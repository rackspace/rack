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

type Environment struct {
    EnvironmentBin []byte
    EnvironmentURL string
    ParsedEnvironment map[string]interface{}
    Files map[string]string
    fileMaps map[string]string
    baseURL string
}

var EnvironmentSections = map[string]bool{
    "parameters": true,
    "parameter_defaults": true,
    "resource_registry": true,
}
/*
Note: although we can define the envirnment format this way, it makes it harder
to make the client have the same behavior as python-heatclient. In this case,
python-heatclient throws an error if the environment contains any keys that are
not the three defined sections
type EnvironmentFormat struct {
    Parameters map[string]string `json:"parameters"`
    ParameterDefaults map[string]string `json:"parameter_defaults"`
    ResourceRegistry map[string]interface{} `json:"resource_registry"`
}
*/
func (e *Environment) Validate() error {
    if e.ParsedEnvironment == nil {
        if err := e.Parse(); err != nil {
            return err
        }
    }
    for key, _ := range e.ParsedEnvironment {
        if _, ok := EnvironmentSections[key]; !ok {
            return errors.New(fmt.Sprintf("Environment has wrong section: %s", key))
        }
    }
    return nil
}

func (e *Environment) Parse() error {
    if err := e.Fetch(); err != nil {
        return err
    }
    if jerr := json.Unmarshal(e.EnvironmentBin, &e.ParsedEnvironment); jerr != nil {
        if yerr := yaml.Unmarshal(e.EnvironmentBin, &e.ParsedEnvironment); yerr != nil {
            return errors.New(fmt.Sprintf("Environment neither json nor yaml format."))
        }
    }
    return e.Validate()
}

func (e *Environment) Fetch() error {
    if e.baseURL == "" {
        basePath, err := filepath.Abs(".")
        if err != nil {
            return err
        }
        u, err := gophercloud.NormalizePathURL("", basePath)
        if err != nil {
            return err
        }
        e.baseURL = u
    }
    u, err := gophercloud.NormalizePathURL(e.baseURL, e.EnvironmentURL)
    if err != nil {
        return err
    }
    e.EnvironmentURL = u

    transport := &http.Transport{}
    transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
    c := &http.Client{Transport: transport}
    resp, err := c.Get(e.EnvironmentURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    e.EnvironmentBin = body
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
    rr := e.ParsedEnvironment["resource_registry"]
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
            tempTemplate := new(Template)
            tempTemplate.baseURL = baseURL
            if err := GetFileContents(tempTemplate, rr, ignoreIfEnvironment); err != nil{
                return err
            }
            if val, ok := rr_map["resources"]; ok {
                switch val.(type){
                case map[string]interface{}, map[interface{}]interface{}:
                    resources_map, err := toStringKeys(val)
                    if err != nil {
                        return err
                    }
                    for _, v := range resources_map {
                        switch v.(type){
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
                            if err := GetFileContents(tempTemplate, v, ignoreIfEnvironment); err != nil{
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

// function to choose keys whose values are other templates
func ignoreIfEnvironment(key string, value interface{}) bool {
	if key == "base_url" || key == "hooks" {
		return true
	}
    if reflect.ValueOf(value).Kind() == reflect.Map {
        return true
    }
    if strValue, ok := value.(string); ok && strings.Contains(strValue, "::") {
        return true
    }
	return false
}

// convert map[interface{}]interface{} to map[string]interface{}
func toStringKeys(m interface{}) (map[string]interface{}, error){
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
func (e *Environment) FixFileRefs() {
    environment_str := string(e.EnvironmentBin)
    for k, v := range e.fileMaps {
        environment_str = strings.Replace(environment_str, k, v, -1)
    }
    e.EnvironmentBin = []byte(environment_str)
}
