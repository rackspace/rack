package stackeventcommands

import (
//    "fmt"
//    "net/http"
	"flag"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
//	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestListResourceContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandListResource{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListResourceKeys(t *testing.T) {
	cmd := &commandListResource{}
	expected := keysListResource
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListResourceServiceClientType(t *testing.T) {
	cmd := &commandListResource{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestListResourceHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
    flagset.String("id", "", "")
    flagset.String("resource", "", "")
	flagset.Set("name", "stack1")
    flagset.Set("id", "id1")
    flagset.Set("resource", "resource1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandListResource{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsListResource{
            stackName: "stack1",
            stackID:    "id1",
            resourceName: "resource1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsListResource{},
	}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsListResource).stackName, actual.Params.(*paramsListResource).stackName)
    th.AssertEquals(t, expected.Params.(*paramsListResource).stackID, actual.Params.(*paramsListResource).stackID)
    th.AssertEquals(t, expected.Params.(*paramsListResource).resourceName, actual.Params.(*paramsListResource).resourceName)
}
/*
func TestListResourceExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1/resources/resource1/events", func(w http.ResponseWriter, r *http.Request) {
            th.TestMethod(t, r, "ListResource")
            w.WriteHeader(http.StatusOK)
            w.Header().Add("Content-Type", "application/json")
            fmt.Fprint(w, `{"events": [{"event_time": "2014-06-03T20:59:46", "resource_name":"resource1"}]}`)
        })
	cmd := &commandListResource{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsListResource{
            stackName: "stack1",
            stackID:    "id1",
            resourceName: "resource1",
		},
	}
	cmd.Execute(actual)
    th.AssertNoErr(t, actual.Err)
}
*/
