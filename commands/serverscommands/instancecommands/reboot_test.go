package instancecommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestRebootContext(t *testing.T) {
	cmd := &commandReboot{
		Ctx: &handler.Context{},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestRebootKeys(t *testing.T) {
	cmd := &commandReboot{}
	expected := keysReboot
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestRebootServiceClientType(t *testing.T) {
	cmd := &commandReboot{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestRebootHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.Bool("soft", false, "")
	flagset.Bool("hard", false, "")
	flagset.Set("soft", "true")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandReboot{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsReboot{
			how: osServers.OSReboot,
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsReboot), *actual.Params.(*paramsReboot))
}

func TestRebootHandlePipe(t *testing.T) {
	cmd := &commandReboot{}
	expected := &handler.Resource{
		Params: &paramsReboot{
			serverID: "server1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsReboot{},
	}
	err := cmd.HandlePipe(actual, "server1")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsReboot).serverID, actual.Params.(*paramsReboot).serverID)
}

func TestRebootHandleSingle(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/servers/detail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"servers":[{"ID":"server1","Name":"server1Name"}]}`)
	})
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "server1Name")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandReboot{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}
	expected := &handler.Resource{
		Params: &paramsReboot{
			serverID: "server1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsReboot{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsReboot).serverID, actual.Params.(*paramsReboot).serverID)
}

func TestRebootExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/servers/server1/action", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	cmd := &commandReboot{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsReboot{
			serverID: "server1",
			how:      "SOFT",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestRebootStdinField(t *testing.T) {
	cmd := &commandReboot{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
