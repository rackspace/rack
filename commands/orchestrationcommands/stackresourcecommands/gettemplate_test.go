package stackresourcecommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
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
	flagset.String("type", "", "")
	flagset.Set("type", "type1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGetTemplate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsGetTemplate{
			resourceType: "type1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGetTemplate{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGetTemplate).resourceType, actual.Params.(*paramsGetTemplate).resourceType)
}

func TestGetTemplateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/resource_types/type1/template", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"HeatTemplateFormatVersion": "2012-12-12"}`)
	})
	cmd := &commandGetTemplate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsGetTemplate{
			resourceType: "type1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestGetTemplateStdinField(t *testing.T) {
	cmd := &commandGetTemplate{}
	expected := "type"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
