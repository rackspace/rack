package containercommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestDeleteContext(t *testing.T) {
	cmd := &commandDelete{
		Ctx: &handler.Context{},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteKeys(t *testing.T) {
	cmd := &commandDelete{}
	expected := keysDelete
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteServiceClientType(t *testing.T) {
	cmd := &commandDelete{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestDeleteHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsDelete{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsDelete), actual.Params.(*paramsDelete))
}

func TestDeleteHandlePipe(t *testing.T) {
	cmd := &commandDelete{}
	expected := &handler.Resource{
		Params: &paramsDelete{
			container: "container1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{},
	}
	err := cmd.HandlePipe(actual, "container1")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).container, actual.Params.(*paramsDelete).container)
}

func TestDeleteHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "container1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsDelete{
			container: "container1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).container, actual.Params.(*paramsDelete).container)
}

func TestDeleteExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/container1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "container1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandDelete{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
			CLIContext:    c,
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{
			container: "container1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestDeleteStdinField(t *testing.T) {
	cmd := &commandDelete{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
