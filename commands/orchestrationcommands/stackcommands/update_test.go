package stackcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestUpdateContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestUpdateKeys(t *testing.T) {
	cmd := &commandUpdate{}
	expected := keysUpdate
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestUpdateServiceClientType(t *testing.T) {
	cmd := &commandUpdate{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestUpdateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			w.WriteHeader(http.StatusAccepted)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, `{"stack": {}}`)
		} else if r.Method == "GET" {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"stack": {"creation_time": "2014-06-03T20:59:46"}}`)
		}
	})
	cmd := &commandUpdate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	templateOpts := new(osStacks.Template)
	templateOpts.Bin = []byte(`"heat_template_version": "2014-10-16"`)
	actual := &handler.Resource{
		Params: &paramsUpdate{
			stackName: "stack1",
			stackID:   "id1",
			opts: &osStacks.UpdateOpts{
				TemplateOpts: templateOpts,
				Parameters: map[string]string{
					"img":    "foo",
					"flavor": "bar",
				},
			},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
