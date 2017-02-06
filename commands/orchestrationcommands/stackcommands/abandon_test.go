package stackcommands

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osStacks "github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestAbandonContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestAbandonKeys(t *testing.T) {
	cmd := &commandAbandon{}
	expected := keysAbandon
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestAbandonServiceClientType(t *testing.T) {
	cmd := &commandAbandon{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestAbandonHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.String("id", "", "")
	flagset.Set("name", "stack1")
	flagset.Set("id", "id1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsAbandon{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsAbandon{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsAbandon).stackName, actual.Params.(*paramsAbandon).stackName)
	th.AssertEquals(t, expected.Params.(*paramsAbandon).stackID, actual.Params.(*paramsAbandon).stackID)
}

func TestAbandonExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1/abandon", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "COMPLETE"}`)
	})

	th.Mux.HandleFunc("/stacks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{
"stacks": [
	{
		"creation_time": "2014-06-03T20:59:46",
		"description": "sample stack",
		"id": "3095aefc-09fb-4bc7-b1f0-f21a304e864c",
		"links": [
			{
				"href": "http://192.168.123.200:8004/v1/eb1c63a4f77141548385f113a28f0f52/stacks/simple_stack/3095aefc-09fb-4bc7-b1f0-f21a304e864c",
				"rel": "self"
			}
		],
		"stack_name": "simple_stack",
		"stack_status": "CREATE_COMPLETE",
		"stack_status_reason": "Stack CREATE completed successfully",
		"updated_time": "",
		"tags": ["foo", "get"]
	}
]
}`)
	})
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsAbandon{
			stackName: "stack1",
			stackID:   "id1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestAbandonStdinField(t *testing.T) {
	cmd := &commandAbandon{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}

func TestAbandonPreTable(t *testing.T) {
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsAbandon{
			stackName: "stack1",
			stackID:   "id1",
		},
	}

	expected := "{\"Action\":\"CREATE\",\"Environment\":{\"encrypted_param_names\":[],\"parameter_defaults\":{},\"parameters\":{},\"resource_registry\":{\"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml\":\"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml\",\"resources\":{}}},\"Files\":{\"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml\":\"heat_template_version: 2014-10-16\\nparameters:\\n  flavor:\\n    type: string\\n    description: Flavor for the server to be created\\n    default: 4353\\n    hidden: true\\nresources:\\n  test_server:\\n    type: \\\"OS::Nova::Server\\\"\\n    properties:\\n      name: test-server\\n      flavor: 2 GB General Purpose v1\\n image: Debian 7 (Wheezy) (PVHVM)\\n\"},\"ID\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Name\":\"postman_stack\",\"ProjectID\":\"897686\",\"Resources\":{\"hello_world\":{\"action\":\"CREATE\",\"name\":\"hello_world\",\"resource_id\":\"8a310d36-46fc-436f-8be4-37a696b8ac63\",\"status\":\"COMPLETE\",\"type\":\"OS::Nova::Server\"}},\"StackUserProjectID\":\"897686\",\"Status\":\"COMPLETE\",\"Template\":{\"description\":\"Simple template to test heat commands\",\"heat_template_version\":\"2013-05-23\",\"parameters\":{\"flavor\":{\"default\":\"m1.tiny\",\"type\":\"string\"}},\"resources\":{\"hello_world\":{\"properties\":{\"flavor\":{\"get_param\":\"flavor\"},\"image\":\"ad091b52-742f-469e-8f3c-fd81cadf0743\",\"key_name\":\"heat_key\",\"user_data\":\"#!/bin/bash -xv\\necho \\\"hello world\\\"; /root/hello-world.txt\\n\"},\"type\":\"OS::Nova::Server\"}}}}"
	resource.Result = &osStacks.AbandonedStack{
		Status: "COMPLETE",
		Name:   "postman_stack",
		Template: map[string]interface{}{
			"heat_template_version": "2013-05-23",
			"description":           "Simple template to test heat commands",
			"parameters": map[string]interface{}{
				"flavor": map[string]interface{}{
					"default": "m1.tiny",
					"type":    "string",
				},
			},
			"resources": map[string]interface{}{
				"hello_world": map[string]interface{}{
					"type": "OS::Nova::Server",
					"properties": map[string]interface{}{
						"key_name": "heat_key",
						"flavor": map[string]interface{}{
							"get_param": "flavor",
						},
						"image":     "ad091b52-742f-469e-8f3c-fd81cadf0743",
						"user_data": "#!/bin/bash -xv\necho \"hello world\"; /root/hello-world.txt\n",
					},
				},
			},
		},
		Action: "CREATE",
		ID:     "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		Resources: map[string]interface{}{
			"hello_world": map[string]interface{}{
				"status":      "COMPLETE",
				"name":        "hello_world",
				"resource_id": "8a310d36-46fc-436f-8be4-37a696b8ac63",
				"action":      "CREATE",
				"type":        "OS::Nova::Server",
			},
		},
		Files: map[string]string{
			"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml": "heat_template_version: 2014-10-16\nparameters:\n  flavor:\n    type: string\n    description: Flavor for the server to be created\n    default: 4353\n    hidden: true\nresources:\n  test_server:\n    type: \"OS::Nova::Server\"\n    properties:\n      name: test-server\n      flavor: 2 GB General Purpose v1\n image: Debian 7 (Wheezy) (PVHVM)\n",
		},
		StackUserProjectID: "897686",
		ProjectID:          "897686",
		Environment: map[string]interface{}{
			"encrypted_param_names": make([]map[string]interface{}, 0),
			"parameter_defaults":    make(map[string]interface{}),
			"parameters":            make(map[string]interface{}),
			"resource_registry": map[string]interface{}{
				"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml": "file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml",
				"resources": make(map[string]interface{}),
			},
		},
	}
	err := cmd.PreTable(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}

func TestAbandonPreJSON(t *testing.T) {
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsAbandon{
			stackName: "stack1",
			stackID:   "id1",
		},
	}

	expected := "{\"action\":\"CREATE\",\"environment\":{\"encrypted_param_names\":[],\"parameter_defaults\":{},\"parameters\":{},\"resource_registry\":{\"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml\":\"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml\",\"resources\":{}}},\"files\":{\"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml\":\"heat_template_version: 2014-10-16\\nparameters:\\n  flavor:\\n    type: string\\n    description: Flavor for the server to be created\\n    default: 4353\\n    hidden: true\\nresources:\\n  test_server:\\n    type: \\\"OS::Nova::Server\\\"\\n    properties:\\n      name: test-server\\n      flavor: 2 GB General Purpose v1\\n image: Debian 7 (Wheezy) (PVHVM)\\n\"},\"id\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"name\":\"postman_stack\",\"project_id\":\"897686\",\"resources\":{\"hello_world\":{\"action\":\"CREATE\",\"name\":\"hello_world\",\"resource_id\":\"8a310d36-46fc-436f-8be4-37a696b8ac63\",\"status\":\"COMPLETE\",\"type\":\"OS::Nova::Server\"}},\"stack_user_project_id\":\"897686\",\"status\":\"COMPLETE\",\"template\":{\"description\":\"Simple template to test heat commands\",\"heat_template_version\":\"2013-05-23\",\"parameters\":{\"flavor\":{\"default\":\"m1.tiny\",\"type\":\"string\"}},\"resources\":{\"hello_world\":{\"properties\":{\"flavor\":{\"get_param\":\"flavor\"},\"image\":\"ad091b52-742f-469e-8f3c-fd81cadf0743\",\"key_name\":\"heat_key\",\"user_data\":\"#!/bin/bash -xv\\necho \\\"hello world\\\" ; /root/hello-world.txt\\n\"},\"type\":\"OS::Nova::Server\"}}}}"
	resource.Result = &osStacks.AbandonedStack{
		Status: "COMPLETE",
		Name:   "postman_stack",
		Template: map[string]interface{}{
			"heat_template_version": "2013-05-23",
			"description":           "Simple template to test heat commands",
			"parameters": map[string]interface{}{
				"flavor": map[string]interface{}{
					"default": "m1.tiny",
					"type":    "string",
				},
			},
			"resources": map[string]interface{}{
				"hello_world": map[string]interface{}{
					"type": "OS::Nova::Server",
					"properties": map[string]interface{}{
						"key_name": "heat_key",
						"flavor": map[string]interface{}{
							"get_param": "flavor",
						},
						"image":     "ad091b52-742f-469e-8f3c-fd81cadf0743",
						"user_data": "#!/bin/bash -xv\necho \"hello world\" ; /root/hello-world.txt\n",
					},
				},
			},
		},
		Action: "CREATE",
		ID:     "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		Resources: map[string]interface{}{
			"hello_world": map[string]interface{}{
				"status":      "COMPLETE",
				"name":        "hello_world",
				"resource_id": "8a310d36-46fc-436f-8be4-37a696b8ac63",
				"action":      "CREATE",
				"type":        "OS::Nova::Server",
			},
		},
		Files: map[string]string{
			"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml": "heat_template_version: 2014-10-16\nparameters:\n  flavor:\n    type: string\n    description: Flavor for the server to be created\n    default: 4353\n    hidden: true\nresources:\n  test_server:\n    type: \"OS::Nova::Server\"\n    properties:\n      name: test-server\n      flavor: 2 GB General Purpose v1\n image: Debian 7 (Wheezy) (PVHVM)\n",
		},
		StackUserProjectID: "897686",
		ProjectID:          "897686",
		Environment: map[string]interface{}{
			"encrypted_param_names": make([]map[string]interface{}, 0),
			"parameter_defaults":    make(map[string]interface{}),
			"parameters":            make(map[string]interface{}),
			"resource_registry": map[string]interface{}{
				"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml": "file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml",
				"resources": make(map[string]interface{}),
			},
		},
	}
	err := cmd.PreJSON(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}

func TestAbandonPreCSV(t *testing.T) {
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsAbandon{
			stackName: "stack1",
			stackID:   "id1",
		},
	}

	expected := "{\"Action\":\"CREATE\",\"Environment\":{\"encrypted_param_names\":[],\"parameter_defaults\":{},\"parameters\":{},\"resource_registry\":{\"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml\":\"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml\",\"resources\":{}}},\"Files\":{\"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml\":\"heat_template_version: 2014-10-16\\nparameters:\\n  flavor:\\n    type: string\\n    description: Flavor for the server to be created\\n    default: 4353\\n    hidden: true\\nresources:\\n  test_server:\\n    type: \\\"OS::Nova::Server\\\"\\n    properties:\\n      name: test-server\\n      flavor: 2 GB General Purpose v1\\n image: Debian 7 (Wheezy) (PVHVM)\\n\"},\"ID\":\"16ef0584-4458-41eb-87c8-0dc8d5f66c87\",\"Name\":\"postman_stack\",\"ProjectID\":\"897686\",\"Resources\":{\"hello_world\":{\"action\":\"CREATE\",\"name\":\"hello_world\",\"resource_id\":\"8a310d36-46fc-436f-8be4-37a696b8ac63\",\"status\":\"COMPLETE\",\"type\":\"OS::Nova::Server\"}},\"StackUserProjectID\":\"897686\",\"Status\":\"COMPLETE\",\"Template\":{\"description\":\"Simple template to test heat commands\",\"heat_template_version\":\"2013-05-23\",\"parameters\":{\"flavor\":{\"default\":\"m1.tiny\",\"type\":\"string\"}},\"resources\":{\"hello_world\":{\"properties\":{\"flavor\":{\"get_param\":\"flavor\"},\"image\":\"ad091b52-742f-469e-8f3c-fd81cadf0743\",\"key_name\":\"heat_key\",\"user_data\":\"#!/bin/bash -xv\\necho \\\"hello world\\\"; /root/hello-world.txt\\n\"},\"type\":\"OS::Nova::Server\"}}}}"
	resource.Result = &osStacks.AbandonedStack{
		Status: "COMPLETE",
		Name:   "postman_stack",
		Template: map[string]interface{}{
			"heat_template_version": "2013-05-23",
			"description":           "Simple template to test heat commands",
			"parameters": map[string]interface{}{
				"flavor": map[string]interface{}{
					"default": "m1.tiny",
					"type":    "string",
				},
			},
			"resources": map[string]interface{}{
				"hello_world": map[string]interface{}{
					"type": "OS::Nova::Server",
					"properties": map[string]interface{}{
						"key_name": "heat_key",
						"flavor": map[string]interface{}{
							"get_param": "flavor",
						},
						"image":     "ad091b52-742f-469e-8f3c-fd81cadf0743",
						"user_data": "#!/bin/bash -xv\necho \"hello world\"; /root/hello-world.txt\n",
					},
				},
			},
		},
		Action: "CREATE",
		ID:     "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
		Resources: map[string]interface{}{
			"hello_world": map[string]interface{}{
				"status":      "COMPLETE",
				"name":        "hello_world",
				"resource_id": "8a310d36-46fc-436f-8be4-37a696b8ac63",
				"action":      "CREATE",
				"type":        "OS::Nova::Server",
			},
		},
		Files: map[string]string{
			"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml": "heat_template_version: 2014-10-16\nparameters:\n  flavor:\n    type: string\n    description: Flavor for the server to be created\n    default: 4353\n    hidden: true\nresources:\n  test_server:\n    type: \"OS::Nova::Server\"\n    properties:\n      name: test-server\n      flavor: 2 GB General Purpose v1\n image: Debian 7 (Wheezy) (PVHVM)\n",
		},
		StackUserProjectID: "897686",
		ProjectID:          "897686",
		Environment: map[string]interface{}{
			"encrypted_param_names": make([]map[string]interface{}, 0),
			"parameter_defaults":    make(map[string]interface{}),
			"parameters":            make(map[string]interface{}),
			"resource_registry": map[string]interface{}{
				"file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml": "file:///Users/prat8228/go/src/github.com/rackspace/rack/my_nova.yaml",
				"resources": make(map[string]interface{}),
			},
		},
	}
	err := cmd.PreCSV(resource)
	th.AssertNoErr(t, err)
	actual, _ := json.Marshal(resource.Result)
	th.AssertEquals(t, expected, string(actual))
}
