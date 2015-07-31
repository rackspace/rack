package objectcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	th "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
	"github.com/jrperritt/rack/output"
)

func newDlCmd(fs *flag.FlagSet) *commandDownload {
	return &commandDownload{Ctx: &handler.Context{
		CLIContext: cli.NewContext(cli.NewApp(), fs, nil),
	}}
}

func TestDlContext(t *testing.T) {
	cmd := newDlCmd(flag.NewFlagSet("flags", 1))
	th.AssertDeepEquals(t, cmd.Ctx, cmd.Context())
}

func TestDlKeys(t *testing.T) {
	cmd := &commandDownload{}
	th.AssertDeepEquals(t, keysDownload, cmd.Keys())
}

func TestDlServiceClientType(t *testing.T) {
	cmd := &commandDownload{}
	th.AssertEquals(t, serviceClientType, cmd.ServiceClientType())
}

func TestDlErrWhenCtnrMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)

	err := newDlCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--container is required."}
	th.AssertDeepEquals(t, expected, err)
}

func TestDlErrWhenNameMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.Set("container", "foo")

	err := newDlCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--name is required."}
	th.AssertDeepEquals(t, expected, err)
}

func TestDlExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "hodor")
	})

	fs := flag.NewFlagSet("flags", 1)
	cmd := newDlCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	actual := &handler.Resource{
		Params: &paramsDownload{container: "foo", object: "bar"},
	}

	cmd.Execute(actual)

	th.AssertNoErr(t, actual.Err)
}
