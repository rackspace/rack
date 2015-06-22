package objectcommands

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osObjects "github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", "[--container <containerName> | stdin container]"),
	Description: "Lists objects in a container",
	Action:      actionList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container containing the objects to list.",
		},
		cli.BoolFlag{
			Name:  "full",
			Usage: "[optional] If false, return just object names. If true, return more complete object information. Defaults to false.",
		},
		cli.StringFlag{
			Name:  "prefix",
			Usage: "[optional] Only return objects with this prefix.",
		},
		cli.StringFlag{
			Name:  "end-marker",
			Usage: "[optional] Only return objects with name less than this value.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing objects at this object name.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many objects at most.",
		},
	}
}

var keysList = []string{"Name", "Bytes", "ContentType", "LastModified"}

type paramsList struct {
	container string
	opts      *osObjects.ListOpts
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
	opts := &osObjects.ListOpts{
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

func (command *commandList) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsList).container = item
	return nil
}

func (command *commandList) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"container"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsList).container = command.Ctx.CLIContext.String("container")
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	containerName := resource.Params.(*paramsList).container
	allPages, err := objects.List(command.Ctx.ServiceClient, containerName, opts).AllPages()
	if err != nil {
		resource.Err = err
		return
	}
	if command.Ctx.CLIContext.IsSet("full") {
		objectInfo, err := objects.ExtractInfo(allPages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(objectInfo))
		for j, obj := range objectInfo {
			result[j] = structs.Map(&obj)
		}
		resource.Result = result
		return
	}
	objectNames, err := objects.ExtractNames(allPages)
	if err != nil {
		resource.Err = err
		return
	}
	result := strings.Join(objectNames, "\n")
	resource.Result = result
}

func (command *commandList) StdinField() string {
	return "container"
}
