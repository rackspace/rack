package stackcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
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

func TestListHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("sort-dir", "", "")
	flagset.String("sort-key", "", "")
	flagset.String("name", "", "")
	flagset.String("status", "", "")
	flagset.String("marker", "", "")
	flagset.Set("sort-dir", "asc")
	flagset.Set("sort-key", "name")
	flagset.Set("name", "stacks*")
	flagset.Set("status", "CREATE_COMPLETE")
	flagset.Set("marker", "1fd3-4f9f-44df-1b5c")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsList{
			opts: &osStacks.ListOpts{
				SortKey: "name",
				SortDir: "asc",
				Name:    "stacks*",
				Status:  "CREATE_COMPLETE",
				Marker:  "1fd3-4f9f-44df-1b5c",
			},
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsList).opts, *actual.Params.(*paramsList).opts)
}

func TestListExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"stacks": [{"creation_time": "2014-06-03T20:59:46"}]}`)
	})
	cmd := &commandList{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsList{
			opts: &osStacks.ListOpts{},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
