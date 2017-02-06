package imagecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--name <imageName> [--server-id <serverID> | --server-name <serverName>]"),
	Description: "Creates an image from a server instance.",
	Action:      actionCreate,
	Flags:       commandoptions.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name the created image will have",
		},
		cli.StringFlag{
			Name:  "server-id",
			Usage: "[optional; required if `server-name` isn't provided] The ID of the server from which to create the image",
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "[optional; required if `server-id` isn't provided] The name of the server from which to create the image",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string of key=value pairs.",
		},
	}
}

var keysCreate = []string{"ID"}

type paramsCreate struct {
	opts     *osServers.CreateImageOpts
	serverID string
}

type commandCreate handler.Command

func actionCreate(c *cli.Context) {
	command := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandCreate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandCreate) Keys() []string {
	return keysCreate
}

func (command *commandCreate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}

	opts := &osServers.CreateImageOpts{
		Name: c.String("name"),
	}

	if c.IsSet("metadata") {
		metadata, err := command.Ctx.CheckKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}

	resource.Params = &paramsCreate{
		opts: opts,
	}

	if c.IsSet("server-id") {
		resource.Params.(*paramsCreate).serverID = c.String("server-id")
		return nil
	}

	serverID, err := osServers.IDFromName(command.Ctx.ServiceClient, c.String("server-name"))
	resource.Params.(*paramsCreate).serverID = serverID

	return err
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	createOpts := resource.Params.(*paramsCreate)
	imageID, err := osServers.CreateImage(command.Ctx.ServiceClient, createOpts.serverID, createOpts.opts).ExtractImageID()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = map[string]interface{}{"ID": imageID}
}
