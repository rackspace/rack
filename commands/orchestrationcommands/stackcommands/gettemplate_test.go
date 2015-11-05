package stackcommands

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

func TestGetTemplateContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGetTemplate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetTemplateKeys(t *testing.T) {
	cmd := &commandGetTemplate{}
	expected := keysGetTemplate
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetTemplateServiceClientType(t *testing.T) {
	cmd := &commandGetTemplate{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestGetTemplateHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.String("id", "", "")
	flagset.Set("name", "stack1")
	flagset.Set("id", "id1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGetTemplate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsGetTemplate{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGetTemplate{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGetTemplate).stackName, actual.Params.(*paramsGetTemplate).stackName)
	th.AssertEquals(t, expected.Params.(*paramsGetTemplate).stackID, actual.Params.(*paramsGetTemplate).stackID)
}

func TestGetTemplateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1/template", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"stack": {"updated_time": "2014-06-03T20:59:46"}}`)
	})
	cmd := &commandGetTemplate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsGetTemplate{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestGetTemplateStdinField(t *testing.T) {
	cmd := &commandGetTemplate{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
