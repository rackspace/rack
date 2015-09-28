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

func TestListContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListKeys(t *testing.T) {
	cmd := &commandList{}
	expected := keysList
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListServiceClientType(t *testing.T) {
	cmd := &commandList{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestListHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("stack-name", "", "")
	flagset.String("stack-id", "", "")
	flagset.Set("stack-name", "stack1")
	flagset.Set("stack-id", "id1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsList{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsList{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsList).stackName, actual.Params.(*paramsList).stackName)
	th.AssertEquals(t, expected.Params.(*paramsList).stackID, actual.Params.(*paramsList).stackID)
}

func TestListExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1/resources", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"resources": [{"creation_time": "2014-06-03T20:59:46", "resource_name":"resource1"}]}`)
	})
	cmd := &commandList{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsList{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestGetStdinField(t *testing.T) {
	cmd := &commandList{}
	expected := "stack-name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
