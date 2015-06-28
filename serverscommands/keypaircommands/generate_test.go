package keypaircommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestGenerateContext(t *testing.T) {
	cmd := &commandGenerate{
		Ctx: &handler.Context{},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGenerateKeys(t *testing.T) {
	cmd := &commandGenerate{}
	expected := keysGenerate
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGenerateServiceClientType(t *testing.T) {
	cmd := &commandGenerate{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestGenerateHandleFlags(t *testing.T) {
	cmd := &commandGenerate{
		Ctx: &handler.Context{},
	}
	expected := &handler.Resource{
		Params: &paramsGenerate{
			opts: &osKeypairs.CreateOpts{},
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsGenerate).opts, *actual.Params.(*paramsGenerate).opts)
}

func TestGenerateHandlePipe(t *testing.T) {
	cmd := &commandGenerate{}
	expected := &handler.Resource{
		Params: &paramsGenerate{
			opts: &osKeypairs.CreateOpts{
				Name: "keypair1Name",
			},
		},
	}
	actual := &handler.Resource{
		Params: &paramsGenerate{
			opts: &osKeypairs.CreateOpts{},
		},
	}
	err := cmd.HandlePipe(actual, "keypair1Name")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, *expected.Params.(*paramsGenerate).opts, *actual.Params.(*paramsGenerate).opts)
}

func TestGenerateHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "keypair1Name")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGenerate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsGenerate{
			opts: &osKeypairs.CreateOpts{
				Name: "keypair1Name",
			},
		},
	}
	actual := &handler.Resource{
		Params: &paramsGenerate{
			opts: &osKeypairs.CreateOpts{},
		},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, *expected.Params.(*paramsGenerate).opts, *actual.Params.(*paramsGenerate).opts)
}

func TestGenerateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/os-keypairs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"keypair":{}}`)
	})
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.Bool("json", false, "")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGenerate{
		Ctx: &handler.Context{
			CLIContext:    c,
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsGenerate{
			opts: &osKeypairs.CreateOpts{
				Name: "keypair1Name",
			},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestGenerateStdinField(t *testing.T) {
	cmd := &commandGenerate{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
