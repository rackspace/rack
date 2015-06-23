package instancecommands

import (
	"flag"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	th "github.com/rackspace/gophercloud/testhelper"
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

func TestGetHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsGet{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsGet), actual.Params.(*paramsGet))
}

func TestGetHandlePipe(t *testing.T) {
	cmd := &commandGet{}
	expected := &handler.Resource{
		Params: &paramsGet{
			server: "server1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{},
	}
	err := cmd.HandlePipe(actual, "server1")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGet).server, actual.Params.(*paramsGet).server)
}

func TestGetHandleSingle(t *testing.T) {
	// need to implement a fake client for HTTP request
	/*
		app := cli.NewApp()
		flagset := flag.NewFlagSet("flags", 1)
		flagset.String("name", "", "")
		flagset.Set("name", "server1")
		c := cli.NewContext(app, flagset, nil)
		cmd := &commandGet{
			Ctx: &handler.Context{
				CLIContext: c,
			},
		}
		expected := &handler.Resource{
			Params: &paramsGet{
				server: "server1",
			},
		}
		actual := &handler.Resource{
			Params: &paramsGet{},
		}
		err := cmd.HandleSingle(actual)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, expected.Params.(*paramsGet).server, actual.Params.(*paramsGet).server)
	*/
}

func TestGetExecute(t *testing.T) {
	// need to implement a fake client for HTTP request
}

func TestGetStdinField(t *testing.T) {
	cmd := &commandGet{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
