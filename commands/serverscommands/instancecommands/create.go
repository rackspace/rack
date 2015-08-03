package instancecommands

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	osServers "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/jrperritt/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "[--name <instanceName> | --stdin name]"),
	Description: "Creates a new server instance",
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
			Usage: "[optional; required if `stdin` isn't provided] The name that the instance should have.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
		cli.StringFlag{
			Name:  "image-id",
			Usage: "[optional; required if `image-name` is not provided] The image ID from which to create the server.",
		},
		cli.StringFlag{
			Name:  "image-name",
			Usage: "[optional; required if `image-id` is not provided] The name of the image from which to create the server.",
		},
		cli.StringFlag{
			Name:  "flavor-id",
			Usage: "[optional; required if `flavor-name` is not provided] The flavor ID that the server should have.",
		},
		cli.StringFlag{
			Name:  "flavor-name",
			Usage: "[optional; required if `flavor-id` is not provided] The name of the flavor that the server should have.",
		},
		cli.StringFlag{
			Name:  "security-groups",
			Usage: "[optional] A comma-separated string of names of the security groups to which this server should belong.",
		},
		cli.StringFlag{
			Name:  "user-data",
			Usage: "[optional] Configuration information or scripts to use after the server boots.",
		},
		cli.StringFlag{
			Name:  "networks",
			Usage: "[optional] A comma-separated string of IDs of the networks to attach to this server. If not provided, a public and private network will be attached.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string of key=value pairs.",
		},
		cli.StringFlag{
			Name:  "admin-pass",
			Usage: "[optional] The root password for the server. If not provided, one will be randomly generated and returned in the output.",
		},
		cli.StringFlag{
			Name:  "keypair",
			Usage: "[optional] The name of the already-existing SSH KeyPair to be injected into this server.",
		},
	}
}

var keysCreate = []string{"ID", "AdminPass"}

type paramsCreate struct {
	opts *servers.CreateOpts
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
	opts := &servers.CreateOpts{
		ImageRef:   c.String("image-id"),
		ImageName:  c.String("image-name"),
		FlavorRef:  c.String("flavor-id"),
		FlavorName: c.String("flavor-name"),
		AdminPass:  c.String("admin-pass"),
		KeyPair:    c.String("keypair"),
	}

	if c.IsSet("security-groups") {
		opts.SecurityGroups = strings.Split(c.String("security-groups"), ",")
	}

	if c.IsSet("user-data") {
		abs, err := filepath.Abs(c.String("user-data"))
		if err != nil {
			return err
		}
		userData, err := ioutil.ReadFile(abs)
		if err != nil {
			return err
		}
		opts.UserData = userData
		opts.ConfigDrive = true
	}

	if c.IsSet("networks") {
		netIDs := strings.Split(c.String("networks"), ",")
		networks := make([]osServers.Network, len(netIDs))
		for i, netID := range netIDs {
			networks[i] = osServers.Network{
				UUID: netID,
			}
		}
		opts.Networks = networks
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
	return nil
}

func (command *commandCreate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsCreate).opts.Name = item
	return nil
}

func (command *commandCreate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsCreate).opts.Name = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	server, err := servers.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		switch err.(type) {
		case *osServers.ErrNeitherImageIDNorImageNameProvided:
			err = errors.New("One and only one of the --image-id and the --image-name flags must be provided.")
		case *osServers.ErrNeitherFlavorIDNorFlavorNameProvided:
			err = errors.New("One and only one of the --flavor-id and the --flavor-name flags must be provided.")
		}
		resource.Err = err
		return
	}
	resource.Result = serverSingle(server)
}

func (command *commandCreate) StdinField() string {
	return "name"
}
