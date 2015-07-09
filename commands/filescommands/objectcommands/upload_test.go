package objectcommands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	th "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
	"github.com/jrperritt/rack/output"
)

func newUpCmd(fs *flag.FlagSet) *commandUpload {
	return &commandUpload{Ctx: &handler.Context{
		CLIContext: cli.NewContext(cli.NewApp(), fs, nil),
	}}
}

func TestUploadContext(t *testing.T) {
	cmd := newUpCmd(flag.NewFlagSet("flags", 1))
	th.AssertDeepEquals(t, cmd.Ctx, cmd.Context())
}

func TestUploadKeys(t *testing.T) {
	cmd := &commandUpload{}
	th.AssertDeepEquals(t, keysUpload, cmd.Keys())
}

func TestUploadServiceClientType(t *testing.T) {
	cmd := &commandUpload{}
	th.AssertEquals(t, serviceClientType, cmd.ServiceClientType())
}

func TestUploadErrWhenCtnrMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)

	err := newUpCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--container is required."}
	th.AssertDeepEquals(t, expected, err)
}

func TestUploadErrWhenNameMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.Set("container", "foo")

	err := newUpCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--name is required."}
	th.AssertDeepEquals(t, expected, err)
}

func TestUploadHandlePipe(t *testing.T) {
	cmd := &commandUpload{}

	actual := &handler.Resource{
		Params: &paramsUpload{},
	}

	err := cmd.HandlePipe(actual, "bar")

	th.AssertNoErr(t, err)
}

func TestUploadHandleSingle(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		th.TestBody(t, r, "body")
		fmt.Fprintf(w, `hodor`)
	})

	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.String("name", "", "")
	fs.String("content", "", "")

	fs.Set("container", "foo")
	fs.Set("name", "bar")
	fs.Set("content", "baz")

	cmd := newUpCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	expected := &handler.Resource{
		Params: &paramsUpload{
			stream: strings.NewReader("baz"),
		},
	}

	actual := &handler.Resource{
		Params: &paramsUpload{},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)

	expectedBytes, _ := ioutil.ReadAll(expected.Params.(*paramsUpload).stream)
	th.AssertNoErr(t, err)

	actualBytes, _ := ioutil.ReadAll(actual.Params.(*paramsUpload).stream)
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, expectedBytes, actualBytes)
}

func TestUploadExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		th.TestMethod(t, r, "PUT")
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, `hodor`)
	})

	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.String("name", "", "")
	fs.Set("container", "foo")
	fs.Set("name", "bar")

	cmd := newUpCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	res := &handler.Resource{
		Params: &paramsUpload{container: "foo", object: "bar"},
	}

	cmd.Execute(res)

	th.AssertNoErr(t, res.Err)
	th.AssertEquals(t, "Successfully uploaded object [bar] to container [foo]\n", res.Result)
}
