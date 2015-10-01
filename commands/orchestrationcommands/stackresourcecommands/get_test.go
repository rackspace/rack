package stackresourcecommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
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
