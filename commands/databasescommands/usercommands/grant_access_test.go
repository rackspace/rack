package usercommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/databases"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestGrantAccessContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGrantAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := cmd.Ctx
	actual := cmd.Context()

	th.AssertDeepEquals(t, expected, actual)
}

func TestGrantAccessKeys(t *testing.T) {
	cmd := &commandGrantAccess{}
	expected := keysGrantAccess
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestGrantAccessServiceClientType(t *testing.T) {
	cmd := &commandGrantAccess{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestGrantAccessHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.String("databases", "", "")

	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")
	flagset.Set("databases", "db1,db2")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGrantAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsGrantAccess{
			userName:   "test",
			instanceID: "instanceId",
			opts: &databases.BatchCreateOpts{
				databases.CreateOpts{Name: "db1"},
				databases.CreateOpts{Name: "db2"},
			},
		},
	}

	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsGrantAccess).opts, *actual.Params.(*paramsGrantAccess).opts)
}

func TestGrantAccessHandlePipe(t *testing.T) {
	cmd := &commandGrantAccess{}
	expected := &handler.Resource{
		Params: &paramsGrantAccess{
			userName: "foo",
		},
	}
	actual := &handler.Resource{
		Params: &paramsGrantAccess{},
	}
	err := cmd.HandlePipe(actual, "foo")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGrantAccess).userName, actual.Params.(*paramsGrantAccess).userName)
}

func TestGrantAccessHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.String("databases", "", "")

	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")
	flagset.Set("databases", "db1,db2")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandGrantAccess{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsGrantAccess{
			userName: "test",
		},
	}

	actual := &handler.Resource{
		Params: &paramsGrantAccess{},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsGrantAccess).userName, actual.Params.(*paramsGrantAccess).userName)
}

func TestGrantAccessExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/users/foo/databases", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
	})

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.String("databases", "", "")
	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")
	flagset.Set("databases", "db1,db2")

	cmd := &commandGrantAccess{
		Ctx: &handler.Context{
			CLIContext:    cli.NewContext(cli.NewApp(), flagset, nil),
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsGrantAccess{
			instanceID: "instanceId",
			userName:   "foo",
			opts: &databases.BatchCreateOpts{
				databases.CreateOpts{Name: "db1"},
			},
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
