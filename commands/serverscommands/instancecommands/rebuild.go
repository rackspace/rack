package instancecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var rebuild = cli.Command{
	Name:        "rebuild",
	Usage:       util.Usage(commandPrefix, "rebuild", "[--id <serverID>|--name <serverName>] --image-id <image-id> --admin-pass <admin-pass>"),
	Description: "Rebuilds an existing server",
	Action:      actionRebuild,
	Flags:       commandoptions.CommandFlags(flagsRebuild, keysGet),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsRebuild, keysGet))
	},
}

func flagsRebuild() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't provided] The ID of the server.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` isn't provided] The name of the server.",
		},
		cli.StringFlag{
			Name:  "image-id",
			Usage: "[required] The ID of the image on which the server will be provisioned.",
		},
		cli.StringFlag{
			Name:  "admin-pass",
			Usage: "[required] The server's administrative password.",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] The name for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "ipv4",
			Usage: "[optional] The IPv4 address for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "ipv6",
			Usage: "[optional] The IPv6 address for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string a key=value pairs.",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the instance has been rebuilt.",
		},
	}
}

var keysRebuild = []string{"ID", "Name", "Status", "Created", "Updated", "Image", "Flavor", "PublicIPv4", "PublicIPv6", "PrivateIPv4", "KeyName"}

type paramsRebuild struct {
	wait     bool
	serverID string
	opts     *servers.RebuildOpts
}

type commandRebuild handler.Command

func actionRebuild(c *cli.Context) {
	command := &commandRebuild{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandRebuild) Context() *handler.Context {
	return command.Ctx
}

func (command *commandRebuild) Keys() []string {
	return keysRebuild
}

func (command *commandRebuild) ServiceClientType() string {
	return serviceClientType
}

func (command *commandRebuild) HandleFlags(resource *handler.Resource) error {
	serverID, err := command.Ctx.IDOrName(osServers.IDFromName)
	if err != nil {
		return err
	}

	err = command.Ctx.CheckFlagsSet([]string{"image-id", "admin-pass"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext
	wait := false
	if c.IsSet("wait-for-completion") {
		wait = true
	}

	opts := &servers.RebuildOpts{
		ImageID:    c.String("image-id"),
		AdminPass:  c.String("admin-pass"),
		AccessIPv4: c.String("ipv4"),
		AccessIPv6: c.String("ipv6"),
	}

	if c.IsSet("metadata") {
		opts.Metadata, err = command.Ctx.CheckKVFlag("metadata")
		if err != nil {
			return err
		}
	}

	if c.IsSet("rename") {
		opts.Name = c.String("rename")
	} else if c.IsSet("name") {
		// Did not set rename, did not set id, can assume name
		opts.Name = c.String("name")
	}

	resource.Params = &paramsRebuild{
		wait:     wait,
		opts:     opts,
		serverID: serverID,
	}

	return nil
}

func (command *commandRebuild) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsRebuild)
	opts := params.opts
	serverID := params.serverID
	server, err := servers.Rebuild(command.Ctx.ServiceClient, serverID, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}

	if params.wait {
		err = osServers.WaitForStatus(command.Ctx.ServiceClient, serverID, "ACTIVE", 600)
		if err != nil {
			resource.Err = err
			return
		}

		server, err = servers.Get(command.Ctx.ServiceClient, serverID).Extract()
		if err != nil {
			resource.Err = err
			return
		}
	}

	resource.Result = serverSingle(server)
}
