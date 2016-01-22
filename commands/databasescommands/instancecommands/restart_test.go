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

func TestRestartContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRestart{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestRestartKeys(t *testing.T) {
	cmd := &commandRestart{}
	expected := keysRestart
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestRestartServiceClientType(t *testing.T) {
	cmd := &commandRestart{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestRestartHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRestart{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsRestart{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsRestart), actual.Params.(*paramsRestart))
}

func TestRestartHandlePipe(t *testing.T) {
	cmd := &commandRestart{}
	expected := &handler.Resource{
		Params: &paramsRestart{
			id: "instanceId",
		},
	}
	actual := &handler.Resource{
		Params: &paramsRestart{},
	}
	err := cmd.HandlePipe(actual, "instanceId")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsRestart).id, actual.Params.(*paramsRestart).id)
}

func TestRestartHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.Set("id", "instanceId")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRestart{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}

	expected := &handler.Resource{
		Params: &paramsRestart{
			id: "instanceId",
		},
	}

	actual := &handler.Resource{
		Params: &paramsRestart{},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsRestart).id, actual.Params.(*paramsRestart).id)
}

func TestRestartExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/action", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	cmd := &commandRestart{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsRestart{
			id: "instanceId",
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestRestartStdinField(t *testing.T) {
	cmd := &commandRestart{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
