package stackresourcecommands

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	osStackResources "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stackresources"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestGetContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetKeys(t *testing.T) {
	cmd := &commandGet{}
	expected := keysGet
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetServiceClientType(t *testing.T) {
	cmd := &commandGet{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestGetHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("stack-name", "", "")
	flagset.String("stack-id", "", "")
	flagset.String("name", "", "")
	flagset.Set("stack-name", "stack1")
	flagset.Set("stack-id", "id1")
	flagset.Set("name", "resource1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsGet{
			stackName:    "stack1",
			stackID:      "id1",
			resourceName: "resource1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{},
	}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGet).stackName, actual.Params.(*paramsGet).stackName)
	th.AssertEquals(t, expected.Params.(*paramsGet).stackID, actual.Params.(*paramsGet).stackID)
	th.AssertEquals(t, expected.Params.(*paramsGet).resourceName, actual.Params.(*paramsGet).resourceName)
}

func TestGetExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1/resources/resource1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"resource": {"updated_time": "2014-06-03T20:59:46"}}`)
	})
	cmd := &commandGet{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{
			stackName:    "stack1",
			stackID:      "id1",
			resourceName: "resource1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestGetPreCSV(t *testing.T) {
	cmd := &commandGet{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsGet{
			stackName:    "stack1",
			stackID:      "id1",
			resourceName: "resource1",
		},
	}

	expected := "{\"Attributes\":{\"SXSW\":\"atx\"},\"CreationTime\":\"\",\"Description\":\"\",\"Links\":[{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance\",\"Rel\":\"self\"},{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e\",\"Rel\":\"stack\"}],\"LogicalID\":\"wordpress_instance\",\"Name\":\"wordpress_instance\",\"PhysicalID\":\"00e3a2fe-c65d-403c-9483-4db9930dd194\",\"RequiredBy\":[],\"Status\":\"CREATE_COMPLETE\",\"StatusReason\":\"state changed\",\"Type\":\"OS::Nova::Server\",\"UpdatedTime\":\"2014-12-10T18:34:35Z\"}"
	resource.Result = &osStackResources.Resource{
		Name: "wordpress_instance",
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance",
				Rel:  "self",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e",
				Rel:  "stack",
			},
		},
		LogicalID:    "wordpress_instance",
		Attributes:   map[string]interface{}{"SXSW": "atx"},
		StatusReason: "state changed",
		UpdatedTime:  time.Date(2014, 12, 10, 18, 34, 35, 0, time.UTC),
		RequiredBy:   []interface{}{},
		Status:       "CREATE_COMPLETE",
		PhysicalID:   "00e3a2fe-c65d-403c-9483-4db9930dd194",
		Type:         "OS::Nova::Server",
	}
	err := cmd.PreCSV(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}

func TestGetPreTable(t *testing.T) {
	cmd := &commandGet{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsGet{},
	}

	expected := "{\"Attributes\":{\"SXSW\":\"atx\"},\"CreationTime\":\"\",\"Description\":\"\",\"Links\":[{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance\",\"Rel\":\"self\"},{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e\",\"Rel\":\"stack\"}],\"LogicalID\":\"wordpress_instance\",\"Name\":\"wordpress_instance\",\"PhysicalID\":\"00e3a2fe-c65d-403c-9483-4db9930dd194\",\"RequiredBy\":[],\"Status\":\"CREATE_COMPLETE\",\"StatusReason\":\"state changed\",\"Type\":\"OS::Nova::Server\",\"UpdatedTime\":\"2014-12-10T18:34:35Z\"}"
	resource.Result = &osStackResources.Resource{
		Name: "wordpress_instance",
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance",
				Rel:  "self",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e",
				Rel:  "stack",
			},
		},
		LogicalID:    "wordpress_instance",
		Attributes:   map[string]interface{}{"SXSW": "atx"},
		StatusReason: "state changed",
		UpdatedTime:  time.Date(2014, 12, 10, 18, 34, 35, 0, time.UTC),
		RequiredBy:   []interface{}{},
		Status:       "CREATE_COMPLETE",
		PhysicalID:   "00e3a2fe-c65d-403c-9483-4db9930dd194",
		Type:         "OS::Nova::Server",
	}

	err := cmd.PreTable(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}
