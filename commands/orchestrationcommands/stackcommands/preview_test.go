package stackcommands

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
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

func TestPreviewPreCSV(t *testing.T) {
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsPreview{},
	}

	expected := "{\"Capabilities\":[],\"CreationTime\":{},\"Description\":\"Simple template to test heat commands\",\"DisableRollback\":true,\"ID\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Links\":[{\"Href\":\"http://166.76.160.117:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Rel\":\"self\"}],\"Name\":\"postman_stack\",\"NotificationTopics\":[],\"Parameters\":{\"OS::stack_id\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"OS::stack_name\":\"postman_stack\",\"flavor\":\"m1.tiny\"},\"Resources\":null,\"TemplateDescription\":\"Simple template to test heat commands\",\"Timeout\":0,\"UpdatedTime\":{}}"
	resource.Result = &osStacks.PreviewedStack{
		DisableRollback: true,
		Description:     "Simple template to test heat commands",
		Parameters: map[string]string{
			"flavor":         "m1.tiny",
			"OS::stack_name": "postman_stack",
			"OS::stack_id":   "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		},
		Name:         "postman_stack",
		CreationTime: time.Date(2015, 2, 3, 20, 7, 39, 0, time.UTC),
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.76.160.117:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/16ef0584-4458-41eb-87c8-0dc8d5f66c87",
				Rel:  "self",
			},
		},
		Capabilities:        []interface{}{},
		NotificationTopics:  []interface{}{},
		ID:                  "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		TemplateDescription: "Simple template to test heat commands",
	}
	err := cmd.PreCSV(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}

func TestPreviewPreTable(t *testing.T) {
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsPreview{},
	}

	expected := "{\"Capabilities\":[],\"CreationTime\":{},\"Description\":\"Simple template to test heat commands\",\"DisableRollback\":true,\"ID\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Links\":[{\"Href\":\"http://166.76.160.117:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Rel\":\"self\"}],\"Name\":\"postman_stack\",\"NotificationTopics\":[],\"Parameters\":{\"OS::stack_id\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"OS::stack_name\":\"postman_stack\",\"flavor\":\"m1.tiny\"},\"Resources\":null,\"TemplateDescription\":\"Simple template to test heat commands\",\"Timeout\":0,\"UpdatedTime\":{}}"
	resource.Result = &osStacks.PreviewedStack{
		DisableRollback: true,
		Description:     "Simple template to test heat commands",
		Parameters: map[string]string{
			"flavor":         "m1.tiny",
			"OS::stack_name": "postman_stack",
			"OS::stack_id":   "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		},
		Name:         "postman_stack",
		CreationTime: time.Date(2015, 2, 3, 20, 7, 39, 0, time.UTC),
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.76.160.117:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/16ef0584-4458-41eb-87c8-0dc8d5f66c87",
				Rel:  "self",
			},
		},
		Capabilities:        []interface{}{},
		NotificationTopics:  []interface{}{},
		ID:                  "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		TemplateDescription: "Simple template to test heat commands",
	}
	err := cmd.PreTable(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}
