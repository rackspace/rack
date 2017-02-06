package snapshotcommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func newDeleteApp(flags map[string]string) *cli.Context {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.String("name", "", "")
	flagset.String("stdin", "", "")
	for k, v := range flags {
		flagset.Set(k, v)
	}
	return cli.NewContext(app, flagset, nil)
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
	c := newDeleteApp(map[string]string{
		"id": "13ba-75c0-4483-acf9",
	})
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
			snapshotID: "13ba-75c0-4483-acf9",
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{},
	}
	err := cmd.HandlePipe(actual, "13ba-75c0-4483-acf9")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).snapshotID, actual.Params.(*paramsDelete).snapshotID)
}

func TestDeleteHandleSingle(t *testing.T) {
	expected := &handler.Resource{
		Params: &paramsDelete{
			snapshotID: "13ba-75c0-4483-acf9",
		},
	}

	actual := &handler.Resource{
		Params: &paramsDelete{},
	}

	cmd := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: newDeleteApp(map[string]string{
				"id": "13ba-75c0-4483-acf9",
			}),
		},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsDelete).snapshotID, actual.Params.(*paramsDelete).snapshotID)
}

func TestDeleteExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/snapshots/13ba-75c0-4483-acf9", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	cmd := &commandDelete{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{
			snapshotID: "13ba-75c0-4483-acf9",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestDeleteStdinField(t *testing.T) {
	cmd := &commandDelete{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
