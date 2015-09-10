package stacks

import (
	"testing"

	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
)

func TestTemplateValidation(t *testing.T) {
	templateJSON := new(Template)
	templateJSON.Bin = []byte(ValidJSONTemplate)
	err := templateJSON.Validate()
	th.AssertNoErr(t, err)

	templateYAML := new(Template)
	templateYAML.Bin = []byte(ValidYAMLTemplate)
	err = templateYAML.Validate()
	th.AssertNoErr(t, err)

	templateInvalid := new(Template)
	templateInvalid.Bin = []byte(InvalidTemplateNoVersion)
	if err = templateInvalid.Validate(); err == nil {
		t.Error("Template validation did not catch invalid template")
	}
}

func TestTemplateParsing(t *testing.T) {
	templateJSON := new(Template)
	templateJSON.Bin = []byte(ValidJSONTemplate)
	err := templateJSON.Parse()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ValidJSONTemplateParsed, templateJSON.Parsed)

	templateYAML := new(Template)
	templateYAML.Bin = []byte(ValidJSONTemplate)
	err = templateYAML.Parse()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ValidJSONTemplateParsed, templateYAML.Parsed)

	templateInvalid := new(Template)
	templateInvalid.Bin = []byte("Keep Austin Weird")
	err = templateInvalid.Parse()
	if err == nil {
		t.Error("Template parsing did not catch invalid template")
	}
}

func TestIgnoreIfTemplate(t *testing.T) {
	var keyValueTests = []struct {
		key   string
		value interface{}
		out   bool
	}{
		{"not_get_file", "afksdf", true},
		{"not_type", "sdfd", true},
		{"get_file", "shdfuisd", false},
		{"type", "dfsdfsd", true},
		{"type", "sdfubsduf.yaml", false},
		{"type", "sdfsdufs.template", false},
		{"type", "sdfsdf.file", true},
		{"type", map[string]string{"key": "value"}, true},
	}
	var result bool
	for _, kv := range keyValueTests {
		result = ignoreIfTemplate(kv.key, kv.value)
		if result != kv.out {
			t.Errorf("key: %v, value: %v expected: %v, actual: %v", kv.key, kv.value, result, kv.out)
		}
	}
}
