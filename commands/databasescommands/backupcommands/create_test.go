package backupcommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/backups"
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

	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.String("description", "", "")

	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")
	flagset.Set("description", "foo")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &backups.CreateOpts{
				Name:        "test",
				InstanceID:  "instanceId",
				Description: "foo",
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
			opts: &backups.CreateOpts{
				Name: "foo",
			},
		},
	}
	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &backups.CreateOpts{},
		},
	}
	err := cmd.HandlePipe(actual, "foo")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsCreate).opts.Name, actual.Params.(*paramsCreate).opts.Name)
}

func TestCreateHandleSingle(t *testing.T) {
	app := cli.NewApp()

	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("instance", "", "")
	flagset.String("name", "", "")
	flagset.String("description", "", "")

	flagset.Set("instance", "instanceId")
	flagset.Set("name", "test")
	flagset.Set("description", "foo")

	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &backups.CreateOpts{
				Name: "test",
			},
		},
	}

	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &backups.CreateOpts{},
		},
	}

	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsCreate).opts.Name, actual.Params.(*paramsCreate).opts.Name)
}

func TestCreateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"backup":{}}`)
	})

	cmd := &commandCreate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}

	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &backups.CreateOpts{
				Name:        "foo",
				InstanceID:  "instanceId",
				Description: "bar",
			},
		},
	}

	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
