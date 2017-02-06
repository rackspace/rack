package keypaircommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestUploadContext(t *testing.T) {
	cmd := &commandUpload{
		Ctx: &handler.Context{},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestUploadKeys(t *testing.T) {
	cmd := &commandUpload{}
	expected := keysUpload
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestUploadServiceClientType(t *testing.T) {
	cmd := &commandUpload{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestUploadHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.String("file", "", "")
	flagset.String("public-key", "", "")
	flagset.Set("name", "keypair1Name")
	flagset.Set("public-key", "ssh public key data here")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandUpload{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsUpload{
			opts: &osKeypairs.CreateOpts{
				Name:      "keypair1Name",
				PublicKey: "ssh public key data here",
			},
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsUpload).opts, *actual.Params.(*paramsUpload).opts)
}

func TestUploadExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/os-keypairs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"keypair":{}}`)
	})
	cmd := &commandUpload{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsUpload{
			opts: &osKeypairs.CreateOpts{
				Name:      "keypair1Name",
				PublicKey: "ssh public key data here",
			},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
