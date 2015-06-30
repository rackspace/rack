package keypaircommands

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jrperritt/rack/handler"
	th "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestListContext(t *testing.T) {
	cmd := &commandList{
		Ctx: &handler.Context{},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListKeys(t *testing.T) {
	cmd := &commandList{}
	expected := keysList
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListServiceClientType(t *testing.T) {
	cmd := &commandList{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestListHandleFlags(t *testing.T) {
}

func TestListHandleSingle(t *testing.T) {
}

func TestListExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/os-keypairs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"keypairs":[]}`)
	})
	cmd := &commandList{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsList{},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
