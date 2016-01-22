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

func TestResizeContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandResize{
		Ctx: &handler.Context{
			CLIContext: c,
		},
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
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandResize{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsResize{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsResize), actual.Params.(*paramsResize))
}

func TestResizeHandlePipe(t *testing.T) {
	cmd := &commandResize{}
	expected := &handler.Resource{
		Params: &paramsResize{
			id: "instanceId",
		},
	}
	actual := &handler.Resource{
		Params: &paramsResize{},
	}
	err := cmd.HandlePipe(actual, "instanceId")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsResize).id, actual.Params.(*paramsResize).id)
}

func TestResizeHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.String("flavor", "", "")
	flagset.Set("id", "instanceId")
	flagset.Set("flavor", "foo")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandResize{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}

	expected := &handler.Resource{
		Params: &paramsResize{
			id:     "instanceId",
			flavor: "foo",
		},
	}

	actual := &handler.Resource{
		Params: &paramsResize{},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsResize).id, actual.Params.(*paramsResize).id)
}

func TestResizeExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/action", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	cmd := &commandResize{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsResize{
			id:     "instanceId",
			flavor: "foo",
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
