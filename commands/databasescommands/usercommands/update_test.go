package usercommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/users"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestUpdateContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := cmd.Ctx
	actual := cmd.Context()

	th.AssertDeepEquals(t, expected, actual)
}

func TestUpdateKeys(t *testing.T) {
	cmd := &commandUpdate{}
	expected := keysUpdate
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestUpdateServiceClientType(t *testing.T) {
	cmd := &commandUpdate{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestUpdateHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.String("new-password", "", "")

	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")
	flagset.Set("new-password", "foo")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsUpdate{
			opts: &users.UpdateOpts{
				Password: "foo",
			},
		},
	}

	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsUpdate).opts, *actual.Params.(*paramsUpdate).opts)
}

func TestUpdateHandlePipe(t *testing.T) {
	cmd := &commandUpdate{}
	expected := &handler.Resource{
		Params: &paramsUpdate{
			opts: &users.UpdateOpts{
				Name: "foo",
			},
		},
	}
	actual := &handler.Resource{
		Params: &paramsUpdate{
			opts: &users.UpdateOpts{},
		},
	}
	err := cmd.HandlePipe(actual, "foo")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsUpdate).opts.Name, actual.Params.(*paramsUpdate).opts.Name)
}

func TestUpdateHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("instance", "", "")
	flagset.String("name", "", "")

	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsUpdate{
			userName: "test",
			opts:     &users.UpdateOpts{},
		},
	}

	actual := &handler.Resource{
		Params: &paramsUpdate{
			userName: "test",
			opts:     &users.UpdateOpts{},
		},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsUpdate).userName, actual.Params.(*paramsUpdate).userName)
}

func TestUpdateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/instanceId/users/name", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
	})

	cmd := &commandUpdate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsUpdate{
			instanceID: "instanceId",
			userName:   "name",
			opts: &users.UpdateOpts{
				Name:     "foo",
				Password: "bar",
			},
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
