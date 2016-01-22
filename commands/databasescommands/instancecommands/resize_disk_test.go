package instancecommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestResizeDiskContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandResizeDisk{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestResizeDiskKeys(t *testing.T) {
	cmd := &commandResizeDisk{}
	expected := keysResizeDisk
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestResizeDiskServiceClientType(t *testing.T) {
	cmd := &commandResizeDisk{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestResizeDiskHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandResizeDisk{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsResizeDisk{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsResizeDisk), actual.Params.(*paramsResizeDisk))
}

func TestResizeDiskHandlePipe(t *testing.T) {
	cmd := &commandResizeDisk{}
	expected := &handler.Resource{
		Params: &paramsResizeDisk{
			id: "instanceId",
		},
	}
	actual := &handler.Resource{
		Params: &paramsResizeDisk{},
	}
	err := cmd.HandlePipe(actual, "instanceId")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsResizeDisk).id, actual.Params.(*paramsResizeDisk).id)
}

func TestResizeDiskHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.Int("size", 0, "")

	flagset.Set("id", "instanceId")
	flagset.Set("size", "5")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandResizeDisk{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}

	expected := &handler.Resource{
		Params: &paramsResizeDisk{
			id: "instanceId",
		},
	}

	actual := &handler.Resource{
		Params: &paramsResizeDisk{},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsResizeDisk).id, actual.Params.(*paramsResizeDisk).id)
}

func TestResizeDiskExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/action", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	cmd := &commandResizeDisk{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsResizeDisk{
			id:   "instanceId",
			size: 5,
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestResizeDiskStdinField(t *testing.T) {
	cmd := &commandResizeDisk{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
