package stackresourcecommands

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

func TestListTypesContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandListTypes{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListTypesKeys(t *testing.T) {
	cmd := &commandListTypes{}
	expected := keysListTypes
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListTypesServiceClientType(t *testing.T) {
	cmd := &commandListTypes{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

/*
func TestListTypesExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1/resources/resource1/events", func(w http.ResponseWriter, r *http.Request) {
            th.TestMethod(t, r, "ListTypes")
            w.WriteHeader(http.StatusOK)
            w.Header().Add("Content-Type", "application/json")
            fmt.Fprint(w, `{"events": [{"event_time": "2014-06-03T20:59:46", "resource_name":"resource1"}]}`)
        })
	cmd := &commandListTypes{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsListTypes{
            stackName: "stack1",
            stackID:    "id1",
            resourceName: "resource1",
		},
	}
	cmd.Execute(actual)
    th.AssertNoErr(t, actual.Err)
}
*/
