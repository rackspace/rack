package snapshotcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/rackspace/gophercloud/rackspace/blockstorage/v1/snapshots"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/rack/handler"
)

func newCreateApp(flags map[string]string) *cli.Context {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("volume-id", "", "")
	flagset.String("name", "", "")
	flagset.String("description", "", "")
	flagset.Bool("force", false, "")
	for k, v := range flags {
		flagset.Set(k, v)
	}
	return cli.NewContext(app, flagset, nil)
}

func TestCreateKeys(t *testing.T) {
	cmd := &commandCreate{}
	expected := keysCreate
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestCreateServiceClientType(t *testing.T) {
	cmd := &commandCreate{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestCreateHandleFlags(t *testing.T) {
	c := newCreateApp(map[string]string{
		"volume-id":   "13ba-75c0-4483-acf9",
		"description": "a description",
		"force":       "true",
	})
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &snapshots.CreateOpts{
				VolumeID:    "13ba-75c0-4483-acf9",
				Description: "a description",
				Force:       true,
			},
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsCreate).opts, *actual.Params.(*paramsCreate).opts)
}

func TestCreateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/snapshots", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"snapshot":{}}`)
	})
	cmd := &commandCreate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &snapshots.CreateOpts{
				VolumeID:    "13ba-75c0-4483-acf9",
				Description: "a description",
			},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
