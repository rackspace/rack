package objectcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/rack/output"
)

func newGetCmd(fs *flag.FlagSet) *commandGet {
	return &commandGet{Ctx: &handler.Context{
		CLIContext: cli.NewContext(cli.NewApp(), fs, nil),
	}}
}

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

func TestGetErrWhenCtnrMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)

	err := newGetCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--container is required."}
	th.AssertDeepEquals(t, expected, err)
}

func TestGetErrWhenNameMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.Set("container", "foo")

	err := newGetCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--name is required."}
	th.AssertDeepEquals(t, expected, err)
}

func TestGetHandlePipe(t *testing.T) {
	cmd := &commandGet{}
	expected := &handler.Resource{
		Params: &paramsGet{object: "bar"},
	}
	actual := &handler.Resource{
		Params: &paramsGet{},
	}

	err := cmd.HandlePipe(actual, "bar")

	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGet).object, actual.Params.(*paramsGet).object)
}

func TestGetHandleSingle(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, `hodor`)
	})

	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.String("name", "", "")
	fs.Set("container", "foo")
	fs.Set("name", "bar")

	cmd := newGetCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	expected := &handler.Resource{
		Params: &paramsGet{
			object: "bar",
		},
	}

	actual := &handler.Resource{
		Params: &paramsGet{},
	}

	err := cmd.HandleSingle(actual)

	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGet).object, actual.Params.(*paramsGet).object)
}

func TestGetExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, `hodor`)
	})

	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.String("name", "", "")
	fs.Set("container", "foo")
	fs.Set("name", "bar")

	cmd := newGetCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	res := &handler.Resource{
		Params: &paramsGet{container: "foo", object: "bar"},
	}

	cmd.Execute(res)

	th.AssertNoErr(t, res.Err)
	th.AssertEquals(t, "text/plain", res.Result.(map[string]interface{})["ContentType"])
}
