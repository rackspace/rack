package backupcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestRestoreContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRestore{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := cmd.Ctx
	actual := cmd.Context()

	th.AssertDeepEquals(t, expected, actual)
}

func TestRestoreKeys(t *testing.T) {
	cmd := &commandRestore{}
	expected := keysRestore
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestRestoreServiceClientType(t *testing.T) {
	cmd := &commandRestore{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestRestoreHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	flagset.String("id", "", "")
	flagset.String("name", "", "")
	flagset.String("flavor", "", "")
	flagset.Int("size", 0, "")

	flagset.Set("id", "id")
	flagset.Set("name", "test")
	flagset.Set("flavor", "flavor-id")
	flagset.Set("size", "5")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRestore{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsRestore{
			opts: &instances.CreateOpts{
				Name:         "test",
				FlavorRef:    "flavor-id",
				Size:         5,
				RestorePoint: "id",
			},
		},
	}

	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsRestore).opts, *actual.Params.(*paramsRestore).opts)
}

func TestRestoreHandlePipe(t *testing.T) {
	cmd := &commandRestore{}
	expected := &handler.Resource{
		Params: &paramsRestore{
			opts: &instances.CreateOpts{
				RestorePoint: "id",
			},
		},
	}
	actual := &handler.Resource{
		Params: &paramsRestore{
			opts: &instances.CreateOpts{},
		},
	}
	err := cmd.HandlePipe(actual, "id")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsRestore).opts.RestorePoint, actual.Params.(*paramsRestore).opts.RestorePoint)
}

func TestRestoreHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.String("name", "", "")
	flagset.String("flavor", "", "")
	flagset.Int("size", 0, "")

	flagset.Set("id", "id")
	flagset.Set("name", "test")
	flagset.Set("flavor", "flavor-id")
	flagset.Set("size", "5")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRestore{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsRestore{
			opts: &instances.CreateOpts{
				RestorePoint: "id",
			},
		},
	}

	actual := &handler.Resource{
		Params: &paramsRestore{
			opts: &instances.CreateOpts{},
		},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsRestore).opts.RestorePoint, actual.Params.(*paramsRestore).opts.RestorePoint)
}

func TestRestoreExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"instance":{}}`)
	})

	cmd := &commandRestore{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsRestore{
			opts: &instances.CreateOpts{
				Name:         "foo",
				FlavorRef:    "bar",
				Size:         10,
				RestorePoint: "id",
			},
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
