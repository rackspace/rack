package stackcommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestDeleteContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteKeys(t *testing.T) {
	cmd := &commandDelete{}
	expected := keysDelete
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteServiceClientType(t *testing.T) {
	cmd := &commandDelete{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestDeleteHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.String("id", "", "")
	flagset.Set("name", "stack1")
	flagset.Set("id", "id1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsDelete{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).stackName, actual.Params.(*paramsDelete).stackName)
	th.AssertEquals(t, expected.Params.(*paramsDelete).stackID, actual.Params.(*paramsDelete).stackID)
}

func TestDeleteExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})

	cmd := &commandDelete{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestDeleteStdinField(t *testing.T) {
	cmd := &commandDelete{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
