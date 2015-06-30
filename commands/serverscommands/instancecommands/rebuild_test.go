package instancecommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/diskconfig"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	th "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestRebuildContext(t *testing.T) {
	cmd := &commandRebuild{
		Ctx: &handler.Context{},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestRebuildKeys(t *testing.T) {
	cmd := &commandRebuild{}
	expected := keysRebuild
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestRebuildServiceClientType(t *testing.T) {
	cmd := &commandRebuild{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestRebuildHandleFlags(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/servers/detail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"servers":[{"ID":"server1","Name":"server1Name"}]}`)
	})

	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("id", "", "")
	flagset.String("image-id", "", "")
	flagset.String("admin-pass", "", "")
	flagset.String("ipv4", "", "")
	flagset.String("ipv6", "", "")
	flagset.String("metadata", "", "")
	flagset.String("rename", "", "")
	flagset.Set("id", "server1")
	flagset.Set("image-id", "123456789")
	flagset.Set("admin-pass", "secret")
	flagset.Set("ipv4", "123.45.67.89")
	flagset.Set("metadata", "img=bar,flavor=foo")
	flagset.Set("rename", "server1Rename")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandRebuild{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsRebuild{
			opts: &servers.RebuildOpts{
				ImageID:    "123456789",
				AdminPass:  "secret",
				AccessIPv4: "123.45.67.89",
				Metadata: map[string]string{
					"img":    "bar",
					"flavor": "foo",
				},
				Name: "server1Rename",
			},
			serverID: "server1",
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsRebuild).opts, *actual.Params.(*paramsRebuild).opts)
}

func TestRebuildExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/servers/server1/action", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"server":{}}`)
	})
	cmd := &commandRebuild{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsRebuild{
			serverID: "server1",
			opts: &servers.RebuildOpts{
				Name:       "server1Rename",
				ImageID:    "123456789",
				AdminPass:  "secret",
				DiskConfig: diskconfig.Auto,
			},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
