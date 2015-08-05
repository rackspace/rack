package objectcommands

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osObjects "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/rack/output"
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

func TestUploadExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		w.Header().Add("Content-Type", "text/plain")
		hash := md5.New()
		io.WriteString(hash, "hodor")
		localChecksum := hash.Sum(nil)
		w.Header().Set("ETag", fmt.Sprintf("%x", localChecksum))
		w.WriteHeader(201)
		fmt.Fprintf(w, `hodor`)
	})

	fs := flag.NewFlagSet("flags", 1)
	cmd := newUpCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	res := &handler.Resource{
		Params: &paramsUpload{
			container: "foo",
			object:    "bar",
			stream:    strings.NewReader("hodor"),
			opts:      osObjects.CreateOpts{},
		},
	}

	cmd.Execute(res)

	th.AssertNoErr(t, res.Err)
	th.AssertEquals(t, "Successfully uploaded object [bar] to container [foo]\n", res.Result)
}
