package usercommands

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

func TestListAccessContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandListAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListAccessKeys(t *testing.T) {
	cmd := &commandListAccess{}
	expected := keysListAccess
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListAccessServiceClientType(t *testing.T) {
	cmd := &commandListAccess{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestListAccessHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.Set("instance", "test")
	flagset.Set("name", "foo")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandListAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsListAccess{
			instanceID: "test",
			userName:   "foo",
		},
	}

	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsListAccess), *actual.Params.(*paramsListAccess))
}

func TestListAccessExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceID/users/foo/databases", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"users":[]}`)
	})

	cmd := &commandListAccess{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsListAccess{
			instanceID: "instanceID",
			userName:   "foo",
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
