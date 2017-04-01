package instancecommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	db "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/databases"
	os "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/instances"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreateContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := cmd.Ctx
	actual := cmd.Context()

	th.AssertDeepEquals(t, expected, actual)
}

func TestCreateKeys(t *testing.T) {
	cmd := &commandCreate{}
	expected := keysCreate
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestCreateServiceClientType(t *testing.T) {
	cmd := &commandCreate{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestCreateHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)

	flagset.String("name", "", "")
	flagset.String("flavor", "", "")
	flagset.Int("size", 0, "")
	flagset.String("config-id", "", "")
	flagset.String("datastore-type", "", "")
	flagset.String("datastore-version", "", "")
	flagset.String("restore-point", "", "")
	flagset.String("replica-of", "", "")

	flagset.Set("name", "test")
	flagset.Set("flavor", "flavor-id")
	flagset.Set("size", "5")
	flagset.Set("config-id", "config-id")
	flagset.Set("datastore-type", "mysql")
	flagset.Set("datastore-version", "5.6")
	flagset.Set("restore-point", "foo")
	flagset.Set("replica-of", "bar")

	flagset.Var(flag.Value(&cli.StringSlice{"db1", "db2", "db3"}), "database", "")
	flagset.Var(flag.Value(&cli.StringSlice{"user1:pw1", "user2:pw2", "user3:pw3"}), "user", "")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &instances.CreateOpts{
				Name:      "test",
				FlavorRef: "flavor-id",
				Size:      5,
				ConfigID:  "config-id",
				Users: users.BatchCreateOpts{
					users.CreateOpts{Name: "user1", Password: "pw1"},
					users.CreateOpts{Name: "user2", Password: "pw2"},
					users.CreateOpts{Name: "user3", Password: "pw3"},
				},
				Databases: db.BatchCreateOpts{
					db.CreateOpts{Name: "db1"},
					db.CreateOpts{Name: "db2"},
					db.CreateOpts{Name: "db3"},
				},
				Datastore: &os.DatastoreOpts{
					Version: "5.6",
					Type:    "mysql",
				},
				RestorePoint: "foo",
				ReplicaOf:    "bar",
			},
		},
	}

	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsCreate).opts, *actual.Params.(*paramsCreate).opts)
}

func TestCreateHandlePipe(t *testing.T) {
	cmd := &commandCreate{}
	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &instances.CreateOpts{
				Name: "instance1",
			},
		},
	}
	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &instances.CreateOpts{},
		},
	}
	err := cmd.HandlePipe(actual, "instance1")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsCreate).opts.Name, actual.Params.(*paramsCreate).opts.Name)
}

func TestCreateHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.String("flavor", "", "")
	flagset.Int("size", 0, "")

	flagset.Set("name", "instance1")
	flagset.Set("flavor", "flavor-id")
	flagset.Set("size", "10")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &instances.CreateOpts{
				Name: "instance1",
			},
		},
	}

	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &instances.CreateOpts{},
		},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsCreate).opts.Name, actual.Params.(*paramsCreate).opts.Name)
}

func TestCreateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"instance":{}}`)
	})

	cmd := &commandCreate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &instances.CreateOpts{
				Name:      "foo",
				FlavorRef: "bar",
				Size:      10,
			},
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
