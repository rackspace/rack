package instancecommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	osServers "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	th "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestResizeContext(t *testing.T) {
	cmd := &commandResize{
		Ctx: &handler.Context{},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestResizeKeys(t *testing.T) {
	cmd := &commandResize{}
	expected := keysResize
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestResizeServiceClientType(t *testing.T) {
	cmd := &commandResize{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestResizeHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("flavor-id", "", "")
	flagset.Set("flavor-id", "2")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandResize{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsResize{
			opts: &osServers.ResizeOpts{
				FlavorRef: "2",
			},
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsResize).opts, *actual.Params.(*paramsResize).opts)
}

func TestResizeHandlePipe(t *testing.T) {
	cmd := &commandResize{}
	expected := &handler.Resource{
		Params: &paramsResize{
			serverID: "server1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsResize{},
	}
	err := cmd.HandlePipe(actual, "server1")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsResize).serverID, actual.Params.(*paramsResize).serverID)
}

func TestResizeHandleSingle(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/servers/detail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"servers":[{"ID":"server1","Name":"server1Name"}]}`)
	})
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.Set("id", "server1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandResize{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}
	expected := &handler.Resource{
		Params: &paramsResize{
			serverID: "server1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsResize{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsResize).serverID, actual.Params.(*paramsResize).serverID)
}

func TestResizeExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/servers/server1/action", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	cmd := &commandResize{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsResize{
			serverID: "server1",
			opts: &osServers.ResizeOpts{
				FlavorRef: "2",
			},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestResizeStdinField(t *testing.T) {
	cmd := &commandResize{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
