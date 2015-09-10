package stacks

import (
	"testing"

	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
)

func TestEnvironmentValidation(t *testing.T) {
	environmentJSON := new(Environment)
	environmentJSON.Bin = []byte(ValidJSONEnvironment)
	err := environmentJSON.Validate()
	th.AssertNoErr(t, err)

	environmentYAML := new(Environment)
	environmentYAML.Bin = []byte(ValidYAMLEnvironment)
	err = environmentYAML.Validate()
	th.AssertNoErr(t, err)

	environmentInvalid := new(Environment)
	environmentInvalid.Bin = []byte(InvalidEnvironment)
	if err = environmentInvalid.Validate(); err == nil {
		t.Error("environment validation did not catch invalid environment")
	}
}

func TestEnvironmentParsing(t *testing.T) {
	environmentJSON := new(Environment)
	environmentJSON.Bin = []byte(ValidJSONEnvironment)
	err := environmentJSON.Parse()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ValidJSONEnvironmentParsed, environmentJSON.Parsed)

	environmentYAML := new(Environment)
	environmentYAML.Bin = []byte(ValidJSONEnvironment)
	err = environmentYAML.Parse()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ValidJSONEnvironmentParsed, environmentYAML.Parsed)

	environmentInvalid := new(Environment)
	environmentInvalid.Bin = []byte("Keep Austin Weird")
	err = environmentInvalid.Parse()
	if err == nil {
		t.Error("environment parsing did not catch invalid environment")
	}
}

func TestIgnoreIfEnvironment(t *testing.T) {
	var keyValueTests = []struct {
		key   string
		value interface{}
		out   bool
	}{
		{"base_url", "afksdf", true},
		{"not_type", "hooks", false},
		{"get_file", "::", true},
		{"hooks", "dfsdfsd", true},
		{"type", "sdfubsduf.yaml", false},
		{"type", "sdfsdufs.environment", false},
		{"type", "sdfsdf.file", false},
		{"type", map[string]string{"key": "value"}, true},
	}
	var result bool
	for _, kv := range keyValueTests {
		result = ignoreIfEnvironment(kv.key, kv.value)
		if result != kv.out {
			t.Errorf("key: %v, value: %v expected: %v, actual: %v", kv.key, kv.value, kv.out, result)
		}
	}
}
