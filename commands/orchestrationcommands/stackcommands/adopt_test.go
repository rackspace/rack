package stackcommands

import (
	"flag"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
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
