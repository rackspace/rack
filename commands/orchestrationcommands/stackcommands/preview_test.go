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

func TestPreviewContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandPreview{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestPreviewKeys(t *testing.T) {
	cmd := &commandPreview{}
	expected := keysPreview
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestPreviewServiceClientType(t *testing.T) {
	cmd := &commandPreview{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestPreviewHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "stack1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandPreview{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsPreview{
			opts: &osStacks.PreviewOpts{Name: "stack1"},
		},
	}
	actual := &handler.Resource{
		Params: &paramsPreview{
			opts: &osStacks.PreviewOpts{},
		},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsPreview).opts.Name, actual.Params.(*paramsPreview).opts.Name)
}

func TestPreviewExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/preview", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"stack": {"id": "3095aefc-09fb-4bc7-b1f0-f21a304e864c"}}`)
	})
	cmd := &commandPreview{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	templateOpts := new(osStacks.Template)
	templateOpts.Bin = []byte(`"heat_template_version": "2014-10-16"`)
	actual := &handler.Resource{
		Params: &paramsPreview{
			opts: &osStacks.PreviewOpts{
				Name:         "stack1",
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

func TestPreviewTemplateStdinField(t *testing.T) {
	cmd := &commandPreview{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
