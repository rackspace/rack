package stackeventcommands

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
	osStackEvents "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stackevents"
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
	flagset.String("resource", "", "")
	flagset.String("id", "", "")
	flagset.Set("stack-name", "stack1")
	flagset.Set("stack-id", "id1")
	flagset.Set("resource", "resource1")
	flagset.Set("id", "event1")
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
			eventID:      "event1",
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
	th.AssertEquals(t, expected.Params.(*paramsGet).eventID, actual.Params.(*paramsGet).eventID)
}

func TestGetExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1/resources/resource1/events/event1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"event": {"event_time": "2014-06-03T20:59:46"}}`)
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
			eventID:      "event1",
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
			eventID:      "event1",
		},
	}

	expected := "{\"ID\":\"93940999-7d40-44ae-8de4-19624e7b8d18\",\"Links\":[{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world/events/93940999-7d40-44ae-8de4-19624e7b8d18\",\"Rel\":\"self\"},{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world\",\"Rel\":\"resource\"},{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b\",\"Rel\":\"stack\"}],\"LogicalResourceID\":\"hello_world\",\"PhysicalResourceID\":\"49181cd6-169a-4130-9455-31185bbfc5bf\",\"ResourceName\":\"hello_world\",\"ResourceProperties\":null,\"ResourceStatus\":\"CREATE_COMPLETE\",\"ResourceStatusReason\":\"state changed\",\"Time\":\"2015-02-05T21:33:27Z\"}"
	resource.Result = &osStackEvents.Event{
		ResourceName: "hello_world",
		Time:         time.Date(2015, 2, 5, 21, 33, 27, 0, time.UTC),
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world/events/93940999-7d40-44ae-8de4-19624e7b8d18",
				Rel:  "self",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
				Rel:  "resource",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
				Rel:  "stack",
			},
		},
		LogicalResourceID:    "hello_world",
		ResourceStatusReason: "state changed",
		ResourceStatus:       "CREATE_COMPLETE",
		PhysicalResourceID:   "49181cd6-169a-4130-9455-31185bbfc5bf",
		ID:                   "93940999-7d40-44ae-8de4-19624e7b8d18",
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

	expected := "{\"ID\":\"93940999-7d40-44ae-8de4-19624e7b8d18\",\"Links\":[{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world/events/93940999-7d40-44ae-8de4-19624e7b8d18\",\"Rel\":\"self\"},{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world\",\"Rel\":\"resource\"},{\"Href\":\"http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b\",\"Rel\":\"stack\"}],\"LogicalResourceID\":\"hello_world\",\"PhysicalResourceID\":\"49181cd6-169a-4130-9455-31185bbfc5bf\",\"ResourceName\":\"hello_world\",\"ResourceProperties\":null,\"ResourceStatus\":\"CREATE_COMPLETE\",\"ResourceStatusReason\":\"state changed\",\"Time\":\"2015-02-05T21:33:27Z\"}"
	resource.Result = &osStackEvents.Event{
		ResourceName: "hello_world",
		Time:         time.Date(2015, 2, 5, 21, 33, 27, 0, time.UTC),
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world/events/93940999-7d40-44ae-8de4-19624e7b8d18",
				Rel:  "self",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
				Rel:  "resource",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
				Rel:  "stack",
			},
		},
		LogicalResourceID:    "hello_world",
		ResourceStatusReason: "state changed",
		ResourceStatus:       "CREATE_COMPLETE",
		PhysicalResourceID:   "49181cd6-169a-4130-9455-31185bbfc5bf",
		ID:                   "93940999-7d40-44ae-8de4-19624e7b8d18",
	}

	err := cmd.PreTable(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}
