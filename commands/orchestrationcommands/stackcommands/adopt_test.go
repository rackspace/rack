package stackcommands

import (
//    "fmt"
//    "net/http"
	"flag"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
//	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestAdoptContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandAdopt{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestAdoptKeys(t *testing.T) {
	cmd := &commandAdopt{}
	expected := keysAdopt
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestAdoptServiceClientType(t *testing.T) {
	cmd := &commandAdopt{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}
/*
func TestAdoptHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
    flagset.String("id", "", "")
	flagset.Set("name", "stack1")
    flagset.Set("id", "id1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandAdopt{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsAdopt{
            stackName: "stack1",
            stackID:    "id1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsAdopt{},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsAdopt).stackName, actual.Params.(*paramsAdopt).stackName)
    th.AssertEquals(t, expected.Params.(*paramsAdopt).stackID, actual.Params.(*paramsAdopt).stackID)
}

func TestAdoptExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1/Adopt", func(w http.ResponseWriter, r *http.Request) {
            th.TestMethod(t, r, "DELETE")
            w.WriteHeader(http.StatusOK)
            w.Header().Add("Content-Type", "application/json")
            fmt.Fprint(w, `{"status": "COMPLETE"}`)
        })
	cmd := &commandAdopt{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsAdopt{
            stackName: "stack1",
            stackID:    "id1",
		},
	}
	cmd.Execute(actual)
    th.AssertNoErr(t, actual.Err)
}
*/
