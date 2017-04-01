package usercommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestRevokeAccessContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRevokeAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := cmd.Ctx
	actual := cmd.Context()

	th.AssertDeepEquals(t, expected, actual)
}

func TestRevokeAccessKeys(t *testing.T) {
	cmd := &commandRevokeAccess{}
	expected := keysRevokeAccess
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestRevokeAccessServiceClientType(t *testing.T) {
	cmd := &commandRevokeAccess{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestRevokeAccessHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.String("database", "", "")

	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")
	flagset.Set("database", "db1")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRevokeAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsRevokeAccess{
			userName:   "test",
			instanceID: "instanceId",
			database:   "db1",
		},
	}

	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsRevokeAccess), *actual.Params.(*paramsRevokeAccess))
}

func TestRevokeAccessHandlePipe(t *testing.T) {
	cmd := &commandRevokeAccess{}
	expected := &handler.Resource{
		Params: &paramsRevokeAccess{
			userName: "foo",
		},
	}
	actual := &handler.Resource{
		Params: &paramsRevokeAccess{},
	}
	err := cmd.HandlePipe(actual, "foo")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsRevokeAccess).userName, actual.Params.(*paramsRevokeAccess).userName)
}

func TestRevokeAccessHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.String("database", "", "")

	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")
	flagset.Set("database", "db1")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRevokeAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsRevokeAccess{
			userName: "test",
		},
	}

	actual := &handler.Resource{
		Params: &paramsRevokeAccess{},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsRevokeAccess).userName, actual.Params.(*paramsRevokeAccess).userName)
}

func TestRevokeAccessExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/users/foo/databases/bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
	})

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.String("databases", "", "")
	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")
	flagset.Set("database", "bar")

	cmd := &commandRevokeAccess{
		Ctx: &handler.Context{
			CLIContext:    cli.NewContext(cli.NewApp(), flagset, nil),
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsRevokeAccess{
			instanceID: "instanceId",
			userName:   "foo",
			database:   "bar",
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
