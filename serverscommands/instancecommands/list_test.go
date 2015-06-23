package instancecommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
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

func TestListHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("image", "", "")
	flagset.String("flavor", "", "")
	flagset.String("name", "", "")
	flagset.String("status", "", "")
	flagset.String("marker", "", "")
	flagset.Set("image", "13ba-75c0-4483-acf9")
	flagset.Set("flavor", "1234-b95f-ac5b-cd23")
	flagset.Set("name", "server*")
	flagset.Set("status", "AVAILABLE")
	flagset.Set("marker", "1fd3-4f9f-44df-1b5c")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsList{
			opts: &osServers.ListOpts{
				Image:  "13ba-75c0-4483-acf9",
				Flavor: "1234-b95f-ac5b-cd23",
				Name:   "server*",
				Status: "AVAILABLE",
				Marker: "1fd3-4f9f-44df-1b5c",
			},
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsList).opts, *actual.Params.(*paramsList).opts)
}

func TestListHandleSingle(t *testing.T) {
	cmd := &commandList{}
	expected := &handler.Resource{}
	actual := &handler.Resource{}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestListExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/servers/detail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{ "servers": [] }`)
	})
	cmd := &commandList{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsList{
			opts: &osServers.ListOpts{},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
