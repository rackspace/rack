package stackcommands

import (
	"flag"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
//	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
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
/*
func TestUpdateHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
    flagset.String("id", "", "")
	flagset.Set("name", "stack1")
    flagset.Set("name", "id1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsUpdate{
            stackName: "stack1",
            stackID: "id1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsUpdate{
            opts: &osStacks.UpdateOpts{},
            stackName: "stack1",
            stackID: "id1",
        },
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsUpdate).stackName, actual.Params.(*paramsUpdate).stackName)
    th.AssertEquals(t, expected.Params.(*paramsUpdate).stackID, actual.Params.(*paramsUpdate).stackID)
}
*/
