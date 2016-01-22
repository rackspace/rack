package databasecommands

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
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
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
			name: "test",
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{},
	}
	err := cmd.HandlePipe(actual, "test")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).name, actual.Params.(*paramsDelete).name)
}

func TestDeleteHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandDelete{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}

	expected := &handler.Resource{
		Params: &paramsDelete{
			instanceID: "instanceId",
			name:       "test",
		},
	}

	actual := &handler.Resource{
		Params: &paramsDelete{},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).name, actual.Params.(*paramsDelete).name)
}

func TestDeleteExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/databases/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	cmd := &commandDelete{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsDelete{
			instanceID: "instanceId",
			name:       "test",
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
