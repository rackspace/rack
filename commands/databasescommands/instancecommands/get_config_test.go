package instancecommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestGetConfigContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGetConfig{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetConfigKeys(t *testing.T) {
	cmd := &commandGetConfig{}
	expected := keysGetConfig
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetConfigServiceClientType(t *testing.T) {
	cmd := &commandGetConfig{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestGetConfigHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGetConfig{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsGetConfig{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsGetConfig), actual.Params.(*paramsGetConfig))
}

func TestGetConfigHandlePipe(t *testing.T) {
	cmd := &commandGetConfig{}
	expected := &handler.Resource{
		Params: &paramsGetConfig{
			id: "instanceId",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGetConfig{},
	}
	err := cmd.HandlePipe(actual, "instanceId")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGetConfig).id, actual.Params.(*paramsGetConfig).id)
}

func TestGetConfigHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.Set("id", "instanceId")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGetConfig{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}
	expected := &handler.Resource{
		Params: &paramsGetConfig{
			id: "instanceId",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGetConfig{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGetConfig).id, actual.Params.(*paramsGetConfig).id)
}

func TestGetConfigExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"instance":{"configuration":{}}}`)
	})

	cmd := &commandGetConfig{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsGetConfig{
			id: "instanceId",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestGetConfigStdinField(t *testing.T) {
	cmd := &commandGetConfig{}
	expected := "id"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
