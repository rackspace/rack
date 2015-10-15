package stackresourcecommands

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStackResources "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stackresources"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestGetSchemaContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGetSchema{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetSchemaKeys(t *testing.T) {
	cmd := &commandGetSchema{}
	expected := keysGetSchema
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetSchemaServiceClientType(t *testing.T) {
	cmd := &commandGetSchema{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestGetSchemaHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("type", "", "")
	flagset.Set("type", "type1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGetSchema{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsGetSchema{
			resourceType: "type1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGetSchema{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGetSchema).resourceType, actual.Params.(*paramsGetSchema).resourceType)
}

func TestGetSchemaExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/resource_types/type1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"resource_type": "OS::Nova::Server"}`)
	})
	cmd := &commandGetSchema{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsGetSchema{
			resourceType: "type1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestGetSchemaStdinField(t *testing.T) {
	cmd := &commandGetSchema{}
	expected := "type"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}

func TestGetSchemaPreCSV(t *testing.T) {
	cmd := &commandGetSchema{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsGetSchema{
			resourceType: "type1",
		},
	}

	expected := "{\"Attributes\":{\"an_attribute\":{\"description\":\"An attribute description .\"}},\"Properties\":{\"a_property\":{\"description\":\"A resource description.\",\"required\":true,\"type\":\"string\",\"update_allowed\":false}},\"ResourceType\":\"OS::Heat::AResourceName\",\"SupportStatus\":{\"message\":\"A status message\",\"status\":\"SUPPORTED\",\"version\":\"2014.1\"}}"
	resource.Result = &osStackResources.TypeSchema{
		Attributes: map[string]interface{}{
			"an_attribute": map[string]interface{}{
				"description": "An attribute description .",
			},
		},
		Properties: map[string]interface{}{
			"a_property": map[string]interface{}{
				"update_allowed": false,
				"required":       true,
				"type":           "string",
				"description":    "A resource description.",
			},
		},
		ResourceType: "OS::Heat::AResourceName",
		SupportStatus: map[string]interface{}{
			"message": "A status message",
			"status":  "SUPPORTED",
			"version": "2014.1",
		},
	}
	err := cmd.PreCSV(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}

func TestGetSchemaPreTable(t *testing.T) {
	cmd := &commandGetSchema{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsGetSchema{
			resourceType: "type1",
		},
	}

	expected := "{\"Attributes\":{\"an_attribute\":{\"description\":\"An attribute description .\"}},\"Properties\":{\"a_property\":{\"description\":\"A resource description.\",\"required\":true,\"type\":\"string\",\"update_allowed\":false}},\"ResourceType\":\"OS::Heat::AResourceName\",\"SupportStatus\":{\"message\":\"A status message\",\"status\":\"SUPPORTED\",\"version\":\"2014.1\"}}"
	resource.Result = &osStackResources.TypeSchema{
		Attributes: map[string]interface{}{
			"an_attribute": map[string]interface{}{
				"description": "An attribute description .",
			},
		},
		Properties: map[string]interface{}{
			"a_property": map[string]interface{}{
				"update_allowed": false,
				"required":       true,
				"type":           "string",
				"description":    "A resource description.",
			},
		},
		ResourceType: "OS::Heat::AResourceName",
		SupportStatus: map[string]interface{}{
			"message": "A status message",
			"status":  "SUPPORTED",
			"version": "2014.1",
		},
	}

	err := cmd.PreTable(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}
