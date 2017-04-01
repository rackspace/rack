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

func TestEnableRootContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandEnableRoot{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestEnableRootKeys(t *testing.T) {
	cmd := &commandEnableRoot{}
	expected := keysEnableRoot
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestEnableRootServiceClientType(t *testing.T) {
	cmd := &commandEnableRoot{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestEnableRootHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandEnableRoot{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsEnableRoot{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsEnableRoot), actual.Params.(*paramsEnableRoot))
}

func TestEnableRootHandlePipe(t *testing.T) {
	cmd := &commandEnableRoot{}
	expected := &handler.Resource{
		Params: &paramsEnableRoot{
			id: "instanceId",
		},
	}
	actual := &handler.Resource{
		Params: &paramsEnableRoot{},
	}
	err := cmd.HandlePipe(actual, "instanceId")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsEnableRoot).id, actual.Params.(*paramsEnableRoot).id)
}

func TestEnableRootHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.Set("id", "instanceId")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandEnableRoot{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}
	expected := &handler.Resource{
		Params: &paramsEnableRoot{
			id: "instanceId",
		},
	}
	actual := &handler.Resource{
		Params: &paramsEnableRoot{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsEnableRoot).id, actual.Params.(*paramsEnableRoot).id)
}

func TestEnableRootExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/root", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"user":{"name":"root","password":"12345"}`)
	})

	cmd := &commandEnableRoot{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsEnableRoot{
			id: "instanceId",
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestEnableRootStdinField(t *testing.T) {
	cmd := &commandEnableRoot{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
