package objectcommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/rack/output"
)

func newDelCmd(fs *flag.FlagSet) *commandDelete {
	return &commandDelete{Ctx: &handler.Context{
		CLIContext: cli.NewContext(cli.NewApp(), fs, nil),
	}}
}

func TestDeleteContext(t *testing.T) {
	cmd := newDelCmd(flag.NewFlagSet("flags", 1))
	th.AssertDeepEquals(t, cmd.Context(), cmd.Ctx)
}

func TestDeleteKeys(t *testing.T) {
	cmd := newDelCmd(flag.NewFlagSet("flags", 1))
	th.AssertDeepEquals(t, cmd.Keys(), keysDelete)
}

func TestDeleteServiceClientType(t *testing.T) {
	cmd := newDelCmd(flag.NewFlagSet("flags", 1))
	th.AssertEquals(t, serviceClientType, cmd.ServiceClientType())
}

func TestDeleteErrWhenCtnrMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)

	err := newDelCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--container is required."}
	th.AssertDeepEquals(t, expected, err)
}

func TestDeleteHandlePipe(t *testing.T) {
	cmd := &commandDelete{}
	expected := &handler.Resource{
		Params: &paramsDelete{object: "bar"},
	}
	actual := &handler.Resource{
		Params: &paramsDelete{},
	}

	err := cmd.HandlePipe(actual, "bar")

	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).object, actual.Params.(*paramsDelete).object)
}

func TestDeleteHandleSingle(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})

	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.String("name", "", "")
	fs.Set("container", "foo")
	fs.Set("name", "bar")

	cmd := newDelCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	expected := &handler.Resource{
		Params: &paramsDelete{
			object: "bar",
		},
	}

	actual := &handler.Resource{
		Params: &paramsDelete{},
	}

	err := cmd.HandleSingle(actual)

	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsDelete).object, actual.Params.(*paramsDelete).object)
}

func TestDeleteExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})

	fs := flag.NewFlagSet("flags", 1)
	fs.String("container", "", "")
	fs.String("name", "", "")
	fs.Set("container", "foo")
	fs.Set("name", "bar")

	cmd := newDelCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	res := &handler.Resource{
		Params: &paramsDelete{container: "foo", object: "bar"},
	}

	cmd.Execute(res)

	th.AssertNoErr(t, res.Err)
}
