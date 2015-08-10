package snapshotcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osSnapshots "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func newListApp(flags map[string]string) *cli.Context {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("volume-id", "", "")
	flagset.String("name", "", "")
	flagset.String("status", "", "")
	for k, v := range flags {
		flagset.Set(k, v)
	}
	return cli.NewContext(app, flagset, nil)
}

func TestListKeys(t *testing.T) {
	cmd := &commandList{}
	expected := keysList
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListServiceClientType(t *testing.T) {
	cmd := &commandList{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestListHandleFlags(t *testing.T) {
	c := newListApp(map[string]string{
		"name":   "rack-test-volume",
		"status": "available",
	})
	cmd := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsList{
			opts: &osSnapshots.ListOpts{
				Name:   "rack-test-volume",
				Status: "available",
			},
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsList).opts, *actual.Params.(*paramsList).opts)
}

func TestListExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/snapshots", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"snapshots": []}`)
	})
	cmd := &commandList{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsList{
			opts: &osSnapshots.ListOpts{
				Name:   "rack-test-volume",
				Status: "available",
			},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
