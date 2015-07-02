package keypaircommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	th "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestDeleteContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteKeys(t *testing.T) {
	cmd := &commandDelete{}
	expected := keysDelete
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteServiceClientType(t *testing.T) {
	cmd := &commandDelete{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestDeleteHandleFlags(t *testing.T) {
	cmd := &commandDelete{
		Ctx: &handler.Context{},
	}
	expected := &handler.Resource{
		Params: &paramsDelete{},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsDelete), actual.Params.(*paramsDelete))
}

func TestDeleteHandlePipe(t *testing.T) {
	cmd := &commandDelete{}
	expected := &handler.Resource{
		Params: &paramsDelete{
			keypair: "keypair1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{},
	}
	err := cmd.HandlePipe(actual, "keypair1")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).keypair, actual.Params.(*paramsDelete).keypair)
}

func TestDeleteHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "keypair1Name")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsDelete{
			keypair: "keypair1Name",
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).keypair, actual.Params.(*paramsDelete).keypair)
}

func TestDeleteExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/os-keypairs/keypair1Name", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	cmd := &commandDelete{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{
			keypair: "keypair1Name",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestDeleteStdinField(t *testing.T) {
	cmd := &commandDelete{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
