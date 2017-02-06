package keypaircommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
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
	cmd := &commandGet{
		Ctx: &handler.Context{},
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
			keypair: "keypair1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{},
	}
	err := cmd.HandlePipe(actual, "keypair1")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGet).keypair, actual.Params.(*paramsGet).keypair)
}

func TestGetHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "keypair1Name")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsGet{
			keypair: "keypair1Name",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGet).keypair, actual.Params.(*paramsGet).keypair)
}

func TestGetExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/os-keypairs/keypair1Name", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"keypair":{}}`)
	})
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGet{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsGet{
			keypair: "keypair1Name",
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
