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
	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.Set("instance", "test")
	flagset.Set("name", "foo")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsGet{
			instanceID: "test",
			userName:   "foo",
		},
	}

	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsGet), *actual.Params.(*paramsGet))
}

func TestGetExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/instances/instanceID/users/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"user":{}}`)
	})

	cmd := &commandGet{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsGet{
			instanceID: "instanceID",
			userName:   "foo",
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
