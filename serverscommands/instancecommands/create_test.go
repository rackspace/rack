package instancecommands

import (
	"flag"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	th "github.com/rackspace/gophercloud/testhelper"
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
	flagset.String("image-id", "", "")
	flagset.String("flavor-id", "", "")
	flagset.String("security-groups", "", "")
	flagset.String("networks", "", "")
	flagset.String("metadata", "", "")
	flagset.String("admin-pass", "", "")
	flagset.String("keypair", "", "")
	flagset.Set("image-id", "foo")
	flagset.Set("flavor-id", "bar")
	flagset.Set("security-groups", "sg1,sg2,sg3")
	flagset.Set("networks", "1111-2222-3333-4444,5555-7777-8888-9999")
	flagset.Set("metadata", "img=foo,flavor=bar")
	flagset.Set("admin-pass", "secret")
	flagset.Set("keypair", "kp1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &servers.CreateOpts{
				ImageRef:       "foo",
				FlavorRef:      "bar",
				SecurityGroups: []string{"sg1", "sg2", "sg3"},
				Networks: []osServers.Network{
					{
						UUID: "1111-2222-3333-4444",
					},
					{
						UUID: "5555-7777-8888-9999",
					},
				},
				Metadata: map[string]string{
					"img":    "foo",
					"flavor": "bar",
				},
				AdminPass: "secret",
				KeyPair:   "kp1",
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
			opts: &servers.CreateOpts{
				Name: "server1",
			},
		},
	}
	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &servers.CreateOpts{},
		},
	}
	err := cmd.HandlePipe(actual, "server1")
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsCreate).opts.Name, actual.Params.(*paramsCreate).opts.Name)
}

func TestCreateHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "server1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &servers.CreateOpts{
				Name: "server1",
			},
		},
	}
	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &servers.CreateOpts{},
		},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsCreate).opts.Name, actual.Params.(*paramsCreate).opts.Name)
}

func TestCreateExecute(t *testing.T) {

}

func StdinField(t *testing.T) {
	cmd := &commandCreate{}
	expected := "name"
	actual := cmd.StdinField()
	th.AssertEquals(t, expected, actual)
}
