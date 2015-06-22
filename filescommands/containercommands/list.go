package containercommands

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osContainers "github.com/rackspace/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists containers",
	Action:      actionList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "full",
			Usage: "If false, return just container names. If true, return more complete container information. Defaults to false.",
		},
		cli.StringFlag{
			Name:  "prefix",
			Usage: "Only return containers with this prefix.",
		},
		cli.StringFlag{
			Name:  "end-marker",
			Usage: "Only return containers with name less than this value.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "Start listing containers at this container name.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "Only return this many containers at most.",
		},
	}
}

var keysList = []string{"Name", "Count", "Bytes"}

type paramsList struct {
	opts *osContainers.ListOpts
}

type commandList handler.Command

func actionList(c *cli.Context) {
	command := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandList) Context() *handler.Context {
	return command.Ctx
}

func (command *commandList) Keys() []string {
	return keysList
}

func (command *commandList) ServiceClientType() string {
	return serviceClientType
}

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	opts := &osContainers.ListOpts{
		Full:      c.Bool("full"),
		Prefix:    c.String("prefix"),
		EndMarker: c.String("end-marker"),
		Marker:    c.String("marker"),
		Limit:     c.Int("limit"),
	}

	resource.Params = &paramsList{
		opts: opts,
	}
	return nil
}

func (command *commandList) HandleSingle(resource *handler.Resource) error {
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	allPages, err := containers.List(command.Ctx.ServiceClient, opts).AllPages()
	if err != nil {
		resource.Err = err
		return
	}
	if command.Ctx.CLIContext.IsSet("full") {
		containers, err := containers.ExtractInfo(allPages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(containers))
		for j, container := range containers {
			result[j] = structs.Map(&container)
		}
		resource.Result = result
	} else {
		containers, err := containers.ExtractNames(allPages)
		if err != nil {
			resource.Err = err
			return
		}
		result := strings.Join(containers, "\n")
		resource.Result = result
	}
}
