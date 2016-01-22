package instancecommands

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

func TestHasRootContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandHasRoot{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestHasRootKeys(t *testing.T) {
	cmd := &commandHasRoot{}
	expected := keysHasRoot
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestHasRootServiceClientType(t *testing.T) {
	cmd := &commandHasRoot{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestHasRootHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandHasRoot{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsHasRoot{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsHasRoot), actual.Params.(*paramsHasRoot))
}

func TestHasRootHandlePipe(t *testing.T) {
	cmd := &commandHasRoot{}
	expected := &handler.Resource{
		Params: &paramsHasRoot{
			id: "instanceId",
		},
	}
	actual := &handler.Resource{
		Params: &paramsHasRoot{},
	}
	err := cmd.HandlePipe(actual, "instanceId")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsHasRoot).id, actual.Params.(*paramsHasRoot).id)
}

func TestHasRootHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.Set("id", "instanceId")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandHasRoot{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}
	expected := &handler.Resource{
		Params: &paramsHasRoot{
			id: "instanceId",
		},
	}
	actual := &handler.Resource{
		Params: &paramsHasRoot{},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsHasRoot).id, actual.Params.(*paramsHasRoot).id)
}

func TestHasRootExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/root", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"rootEnabled": true}`)
	})

	cmd := &commandHasRoot{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsHasRoot{
			id: "instanceId",
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestHasRootStdinField(t *testing.T) {
	cmd := &commandHasRoot{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
