package instancecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description: "Retrieves an existing server",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the server.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the server.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
	}
}

var keysGet = []string{"ID", "Name", "Status", "Created", "Updated", "Image", "Flavor",
	"PublicIPv4", "PublicIPv6", "PrivateIPv4", "KeyName", "TenantID", "Progress",
	"SecurityGroups", "Metadata"}

type paramsGet struct {
	server string
}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGet).server = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	id, err := command.Context().IDOrName(osServers.IDFromName)
	resource.Params.(*paramsGet).server = id
	return err
}

func (command *commandGet) Execute(resource *handler.Resource) {
	serverID := resource.Params.(*paramsGet).server
	server, err := servers.Get(command.Context().ServiceClient, serverID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = serverSingle(server)
}

func (command *commandGet) StdinField() string {
	return "id"
}

func (command *commandGet) PreCSV(resource *handler.Resource) error {
	resource.FlattenMap("Metadata")
	return nil
}

func (command *commandGet) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
