package snapshotcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func newGetApp(flags map[string]string) *cli.Context {
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

func TestGetKeys(t *testing.T) {
	cmd := &commandGet{}
	expected := keysGet
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetServiceClientType(t *testing.T) {
	cmd := &commandGet{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestGetHandleFlags(t *testing.T) {
	c := newGetApp(map[string]string{
		"id": "13ba-75c0-4483-acf9",
	})
	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsGet{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsGet), actual.Params.(*paramsGet))
}

func TestGetHandlePipe(t *testing.T) {
	cmd := &commandGet{}
	expected := &handler.Resource{
		Params: &paramsGet{
			snapshotID: "13ba-75c0-4483-acf9",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{},
	}
	err := cmd.HandlePipe(actual, "13ba-75c0-4483-acf9")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGet).snapshotID, actual.Params.(*paramsGet).snapshotID)
}

func TestGetHandleSingle(t *testing.T) {
	expected := &handler.Resource{
		Params: &paramsGet{
			snapshotID: "13ba-75c0-4483-acf9",
		},
	}

	actual := &handler.Resource{
		Params: &paramsGet{},
	}

	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext: newGetApp(map[string]string{
				"id": "13ba-75c0-4483-acf9",
			}),
		},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsGet).snapshotID, actual.Params.(*paramsGet).snapshotID)
}

func TestGetExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/snapshots/13ba-75c0-4483-acf9", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"snapshot":{}}`)
	})
	cmd := &commandGet{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{
			snapshotID: "13ba-75c0-4483-acf9",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestGetStdinField(t *testing.T) {
	cmd := &commandGet{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
