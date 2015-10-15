package stackcommands

import (
	"encoding/json"
	"flag"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestAdoptContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandAdopt{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestAdoptKeys(t *testing.T) {
	cmd := &commandAdopt{}
	expected := keysAdopt
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestAdoptServiceClientType(t *testing.T) {
	cmd := &commandAdopt{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestAdoptPreCSV(t *testing.T) {
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsAdopt{},
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

func TestAdoptPreTable(t *testing.T) {
	cmd := &commandAbandon{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	resource := &handler.Resource{
		Params: &paramsAdopt{},
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
