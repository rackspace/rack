package instancecommands

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var rebuild = cli.Command{
	Name:        "rebuild",
	Usage:       util.Usage(commandPrefix, "rebuild", strings.Join([]string{util.IDOrNameUsage("instance"), "--image-id <image-id>", "--admin-pass <admin-pass>"}, " ")),
	Description: "Rebuilds an existing server",
	Action:      actionRebuild,
	Flags:       util.CommandFlags(flagsRebuild, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsRebuild, keysGet))
	},
}

func flagsRebuild() []cli.Flag {
	cf := []cli.Flag{
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
	}
	return append(cf, util.IDAndNameFlags...)
}

var keysRebuild = []string{"ID", "Name", "Status", "Created", "Updated", "Image", "Flavor", "Public IPv4", "Public IPv6", "Private IPv4", "KeyName"}

type paramsRebuild struct {
	serverID string
	opts     *osServers.RebuildOpts
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
	c := command.Context().CLIContext
	err := command.Context().CheckFlagsSet([]string{"image-id", "admin-pass"})
	if err != nil {
		return err
	}
	opts := &osServers.RebuildOpts{
		ImageID:    c.String("image-id"),
		AdminPass:  c.String("admin-pass"),
		AccessIPv4: c.String("ipv4"),
		AccessIPv6: c.String("ipv6"),
	}
	if c.IsSet("metadata") {
		opts.Metadata, err = command.Context().CheckKVFlag("metadata")
		if err != nil {
			return err
		}
	}
	resource.Params = &paramsRebuild{
		opts: opts,
	}
	return nil
}

func (command *commandRebuild) HandleSingle(resource *handler.Resource) error {
	id, err := command.Context().IDOrName(osServers.IDFromName)
	resource.Params.(*paramsRebuild).serverID = id
	return err
}

func (command *commandRebuild) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsRebuild).opts
	serverID := resource.Params.(*paramsRebuild).serverID
	c := command.Context().CLIContext
	if c.IsSet("rename") {
		opts.Name = c.String("rename")
	} else if c.IsSet("id") { // Must get the name from compute by ID
		server, err := servers.Get(command.Context().ServiceClient, serverID).Extract()
		if err != nil {
			resource.Err = err
			return
		}
		opts.Name = server.Name
	} else if c.IsSet("name") {
		// Did not set rename, did not set id, can assume name
		opts.Name = c.String("name")
	}
	server, err := servers.Rebuild(command.Context().ServiceClient, serverID, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = serverSingle(server)
	command.Context().Results <- resource
}
