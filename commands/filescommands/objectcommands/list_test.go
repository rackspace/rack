package objectcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osObjects "github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/rack/output"
)

func setupList(t *testing.T) {
	th.Mux.HandleFunc("/testContainer", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `[{"name": "goodbye"},{"name": "hello"}]`)
		case "hello":
			fmt.Fprintf(w, `[]`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

func newListCmd(fs *flag.FlagSet) *commandList {
	return &commandList{Ctx: &handler.Context{
		CLIContext: cli.NewContext(cli.NewApp(), fs, nil),
	}}
}

func TestListContext(t *testing.T) {
	cmd := newListCmd(flag.NewFlagSet("flags", 1))
	th.AssertDeepEquals(t, cmd.Ctx, cmd.Context())
}

func TestListKeys(t *testing.T) {
	cmd := &commandList{}
	th.AssertDeepEquals(t, keysList, cmd.Keys())
}

func TestListServiceClientType(t *testing.T) {
	cmd := &commandList{}
	th.AssertEquals(t, serviceClientType, cmd.ServiceClientType())
}

func TestListErrWhenCtnrMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)

	err := newListCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--container is required."}
	th.AssertDeepEquals(t, expected, err)
}

func TestListHandlePipe(t *testing.T) {
	cmd := &commandList{}
	expected := &handler.Resource{
		Params: &paramsList{container: "bar"},
	}
	actual := &handler.Resource{
		Params: &paramsList{},
	}

	err := cmd.HandlePipe(actual, "bar")

	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsList).container, actual.Params.(*paramsList).container)
}

func TestListHandleSingle(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	setupList(t)

	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.Set("container", "testContainer")

	cmd := newListCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	expected := &handler.Resource{
		Params: &paramsList{
			container: "testContainer",
		},
	}

	actual := &handler.Resource{
		Params: &paramsList{},
	}

	err := cmd.HandleSingle(actual)

	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsList).container, actual.Params.(*paramsList).container)
}

func TestListExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	setupList(t)

	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.Bool("all-pages", false, "")
	fs.Set("container", "foo")
	fs.Set("all-pages", "true")

	cmd := newListCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	res := &handler.Resource{
		Params: &paramsList{container: "testContainer", opts: &osObjects.ListOpts{}},
	}

	cmd.Execute(res)

	th.AssertNoErr(t, res.Err)
	th.AssertEquals(t, 2, len(res.Result.([]map[string]interface{})))
}
