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

var TemplateFormatVersions = map[string]bool{
    "HeatTemplateFormatVersion": true,
    "heat_template_version": true,
    "AWSTemplateFormatVersion": true,
}

func (t *Template) Validate() error {
    if t.ParsedTemplate == nil {
        if err := t.Parse(); err != nil {
            return err
        }
    }
    for key, _ := range t.ParsedTemplate {
        if _, ok := TemplateFormatVersions[key]; ok {
            return nil
        }
    }
    return errors.New(fmt.Sprintf("Template format version not found."))
}

type Template struct {
    TemplateBin []byte
    TemplateURL string
    ParsedTemplate map[string]interface{}  // internal representation of template
    Files map[string]string
    fileMaps map[string]string
    baseURL string
}

func (t *Template) Fetch() error {
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
    u, err := gophercloud.NormalizePathURL(t.baseURL, t.TemplateURL)
    if err != nil {
        return err
    }
    t.TemplateURL = u

    transport := &http.Transport{}
    transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
    c := &http.Client{Transport: transport}
    resp, err := c.Get(t.TemplateURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    t.TemplateBin = body
    return nil
}

func (t *Template) Parse() error {
    if err := t.Fetch(); err != nil {
        return err
    }
    if jerr := json.Unmarshal(t.TemplateBin, &t.ParsedTemplate); jerr != nil {
        if yerr := yaml.Unmarshal(t.TemplateBin, &t.ParsedTemplate); yerr != nil {
            return errors.New(fmt.Sprintf("Template neither json nor yaml format."))
        }
    }
    return t.Validate()
}

type igFunc func(string, interface{}) bool

func GetFileContents(t *Template, template interface{}, ignoreIf igFunc) error {
    if t.Files == nil {
        t.Files = make(map[string]string)
    }
    if t.fileMaps == nil {
        t.fileMaps = make(map[string]string)
    }
    switch template.(type) {
    case map[string]interface{}, map[interface{}]interface{}:
            template_map, err := toStringKeys(template)
            if err != nil {
                return err
            }

            for k, v := range template_map {
                value, ok := v.(string)
                if !ok {
                    if err := GetFileContents(t, v, ignoreIf); err != nil {
                        return err
                    }
                } else if !ignoreIf(k, value) {
                    // at this point, the k, v pair has a reference to an external template.
                    // The assumption of heatclient is that value v is a relative reference
                    // to a file in the users environment
                    childtemplate := new(Template)
                    baseURL, err :=  gophercloud.NormalizePathURL(t.baseURL, value)
                    if err != nil {
                        return err
                    }
                    childtemplate.baseURL = baseURL
                    if err := childtemplate.Parse(); err != nil {
                        //TODO: log this error!
                        return nil
                    }
                    if err := GetFileContents(childtemplate, childtemplate.ParsedTemplate, ignoreIfTemplate); err != nil {
                        return err
                    }
                    // update parent template with current child templates' content
                    t.fileMaps[value] = childtemplate.TemplateURL
                    t.Files[childtemplate.TemplateURL] = string(childtemplate.TemplateBin)
                }
            }
            return nil
        case []interface{}:
            template_slice := template.([]interface{})
            for i := range template_slice {
                if err := GetFileContents(t, template_slice[i], ignoreIfTemplate); err != nil {
                    return err
                }
            }
        case string, bool, float64, nil, int:
            return nil
        default:
            return errors.New(fmt.Sprintf("%v: Unrecognized type", reflect.TypeOf(template)))

    }
    return nil
}

// fix the template reference to files
func (t *Template) FixFileRefs() {
    template_str := string(t.TemplateBin)
    for k, v := range t.fileMaps {
        template_str = strings.Replace(template_str, k, v, -1)
    }
    t.TemplateBin = []byte(template_str)
}

// function to choose keys whose values are other templates
func ignoreIfTemplate(key string, value interface{}) bool {
	if key != "get_file" && key != "type" {
		return true
	}
	if key == "type" && !(strings.HasSuffix(value.(string), ".template") || strings.HasSuffix(value.(string), ".yaml")) {
		return true
	}
	return false
}
