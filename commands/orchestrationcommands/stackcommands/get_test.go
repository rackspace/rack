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

func TestGetContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
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

func TestGetHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.String("id", "", "")
	flagset.Set("name", "stack1")
	flagset.Set("id", "id1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsGet{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGet).stackName, actual.Params.(*paramsGet).stackName)
	th.AssertEquals(t, expected.Params.(*paramsGet).stackID, actual.Params.(*paramsGet).stackID)
}

func TestGetExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"stack": {"creation_time": "2014-06-03T20:59:46"}}`)
	})
	cmd := &commandGet{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestGetStdinField(t *testing.T) {
	cmd := &commandGet{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}

func TestGetPreCSV(t *testing.T) {
	cmd := &commandGet{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsGet{
			stackName: "stack1",
			stackID:   "id1",
		},
	}

	expected := "{\"Capabilities\":[],\"CreationTime\":\"2015-02-03T20:07:39Z\",\"Description\":\"Simple template to test heat commands\",\"DisableRollback\":true,\"ID\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Links\":[{\"Href\":\"http://166.76.160.117:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Rel\":\"self\"}],\"Name\":\"postman_stack\",\"NotificationTopics\":[],\"Outputs\":[],\"Parameters\":{\"OS::stack_id\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"OS::stack_name\":\"postman_stack\",\"flavor\":\"m1.tiny\"},\"Status\":\"CREATE_COMPLETE\",\"StatusReason\":\"Stack Get completed successfully\",\"Tags\":[\"rackspace\",\"atx\"],\"TemplateDescription\":\"Simple template to test heat commands\",\"Timeout\":0,\"UpdatedTime\":\"\"}"
	resource.Result = &osStacks.RetrievedStack{
		DisableRollback: true,
		Description:     "Simple template to test heat commands",
		Parameters: map[string]string{
			"flavor":         "m1.tiny",
			"OS::stack_name": "postman_stack",
			"OS::stack_id":   "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		},
		StatusReason: "Stack Get completed successfully",
		Name:         "postman_stack",
		Outputs:      []map[string]interface{}{},
		CreationTime: time.Date(2015, 2, 3, 20, 7, 39, 0, time.UTC),
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.76.160.117:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/16ef0584-4458-41eb-87c8-0dc8d5f66c87",
				Rel:  "self",
			},
		},
		Capabilities:        []interface{}{},
		NotificationTopics:  []interface{}{},
		Status:              "CREATE_COMPLETE",
		ID:                  "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		TemplateDescription: "Simple template to test heat commands",
		Tags:                []string{"rackspace", "atx"},
	}
	err := cmd.PreCSV(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}

func TestGetPreTable(t *testing.T) {
	cmd := &commandGet{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsGet{
			stackName: "stack1",
			stackID:   "id1",
		},
	}

	expected := "{\"Capabilities\":[],\"CreationTime\":\"2015-02-03T20:07:39Z\",\"Description\":\"Simple template to test heat commands\",\"DisableRollback\":true,\"ID\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Links\":[{\"Href\":\"http://166.76.160.117:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Rel\":\"self\"}],\"Name\":\"postman_stack\",\"NotificationTopics\":[],\"Outputs\":[],\"Parameters\":{\"OS::stack_id\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"OS::stack_name\":\"postman_stack\",\"flavor\":\"m1.tiny\"},\"Status\":\"CREATE_COMPLETE\",\"StatusReason\":\"Stack CREATE completed successfully\",\"Tags\":[\"rackspace\",\"atx\"],\"TemplateDescription\":\"Simple template to test heat commands\",\"Timeout\":0,\"UpdatedTime\":\"\"}"
	resource.Result = &osStacks.RetrievedStack{
		DisableRollback: true,
		Description:     "Simple template to test heat commands",
		Parameters: map[string]string{
			"flavor":         "m1.tiny",
			"OS::stack_name": "postman_stack",
			"OS::stack_id":   "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		},
		StatusReason: "Stack CREATE completed successfully",
		Name:         "postman_stack",
		Outputs:      []map[string]interface{}{},
		CreationTime: time.Date(2015, 2, 3, 20, 7, 39, 0, time.UTC),
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.76.160.117:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/16ef0584-4458-41eb-87c8-0dc8d5f66c87",
				Rel:  "self",
			},
		},
		Capabilities:        []interface{}{},
		NotificationTopics:  []interface{}{},
		Status:              "CREATE_COMPLETE",
		ID:                  "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		TemplateDescription: "Simple template to test heat commands",
		Tags:                []string{"rackspace", "atx"},
	}

	err := cmd.PreTable(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}
