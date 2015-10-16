package stackcommands

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreateContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
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
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.String("disable-rollback", "", "")
	flagset.String("template-file", "", "")
	flagset.String("environment-file", "", "")
	flagset.String("timeout", "", "")
	flagset.String("parameters", "", "")
	flagset.Set("name", "stack1")
	flagset.Set("disable-rollback", "true")
	flagset.Set("template-file", "mytemplate.yaml")
	flagset.Set("environment-file", "myenvironment.yaml")
	flagset.Set("timeout", "300")
	flagset.Set("parameters", "img=foo,flavor=bar")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	disableRollback := true
	templateOpts := new(osStacks.Template)
	templateOpts.URL = "mytemplate.yaml"
	environmentOpts := new(osStacks.Environment)
	environmentOpts.URL = "myenvironment.yaml"
	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &osStacks.CreateOpts{
				Name:            "stack1",
				TemplateOpts:    templateOpts,
				EnvironmentOpts: environmentOpts,
				DisableRollback: &disableRollback,
				Timeout:         300,
				Parameters: map[string]string{
					"img":    "foo",
					"flavor": "bar",
				},
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
	th.Mux.HandleFunc("/stacks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, `{"stack": {"id": "3095aefc-09fb-4bc7-b1f0-f21a304e864c"}}`)
	})
	cmd := &commandCreate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	templateOpts := new(osStacks.Template)
	templateOpts.Bin = []byte(`"heat_template_version": "2014-10-16"`)
	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &osStacks.CreateOpts{
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

func TestCreatePreCSV(t *testing.T) {
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsCreate{},
	}

	expected := "{\"ID\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Links\":[{\"Href\":\"http://168.28.170.117:8004/v1/98606384f58drad0bhdb7d02779549ac/stacks/stackcreated/16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Rel\":\"self\"}]}"
	resource.Result = &osStacks.CreatedStack{
		ID: "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://168.28.170.117:8004/v1/98606384f58drad0bhdb7d02779549ac/stacks/stackcreated/16ef0584-4458-41eb-87c8-0dc8d5f66c87",
				Rel:  "self",
			},
		},
	}
	err := cmd.PreCSV(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}

func TestCreatePreTable(t *testing.T) {
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsCreate{},
	}

	expected := "{\"ID\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Links\":[{\"Href\":\"http://168.28.170.117:8004/v1/98606384f58drad0bhdb7d02779549ac/stacks/stackcreated/16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Rel\":\"self\"}]}"
	resource.Result = &osStacks.CreatedStack{
		ID: "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://168.28.170.117:8004/v1/98606384f58drad0bhdb7d02779549ac/stacks/stackcreated/16ef0584-4458-41eb-87c8-0dc8d5f66c87",
				Rel:  "self",
			},
		},
	}
	err := cmd.PreTable(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}
