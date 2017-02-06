package containercommands

import (
	"flag"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreateContext(t *testing.T) {
	cmd := &commandCreate{
		Ctx: &handler.Context{},
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
	flagset.String("metadata", "", "")
	flagset.String("container-read", "", "")
	flagset.String("container-write", "", "")
	flagset.Set("metadata", "key=val,foo=bar")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: containers.CreateOpts{
				Metadata: map[string]string{
					"key": "val",
					"foo": "bar",
				},
			},
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected.Params.(*paramsCreate).opts, actual.Params.(*paramsCreate).opts)
}

func TestCreateHandlePipe(t *testing.T) {
	cmd := &commandCreate{}
	expected := &handler.Resource{
		Params: &paramsCreate{
			container: "container1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsCreate{},
	}
	err := cmd.HandlePipe(actual, "container1")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsCreate).container, actual.Params.(*paramsCreate).container)
}

func TestCreateHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "container1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsCreate{
			container: "container1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsCreate{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsCreate).container, actual.Params.(*paramsCreate).container)
}

func TestCreateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/container1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Container-Meta-Foo", "bar")
		w.Header().Add("X-Container-Meta-Key", "val")
		w.Header().Add("X-Trans-Id", "1234567")
		w.WriteHeader(http.StatusNoContent)
	})
	cmd := &commandCreate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: containers.CreateOpts{
				Metadata: map[string]string{
					"key": "val",
					"foo": "bar",
				},
			},
			container: "container1",
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}

func TestCreateStdinField(t *testing.T) {
	cmd := &commandCreate{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
