package instancecommands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osBFV "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	osServers "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	bfv "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/bootfromvolume"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/rackspace/rack/util"
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
			Usage: "[optional; required if `image-name` or `block-device` is not provided] The image ID from which to create the server.",
		},
		cli.StringFlag{
			Name:  "image-name",
			Usage: "[optional; required if `image-id` or `block-device` is not provided] The name of the image from which to create the server.",
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
		cli.StringFlag{
			Name: "block-device",
			Usage: strings.Join([]string{"[optional] Used to boot from volume.",
				"\tIf provided, the instance will be created based upon the comma-separated key=value pairs provided to this flag.",
				"\tOptions:",
				"\t\tsource-type\t[required] The source type of the device. Options: volume, snapshot, image.",
				"\t\tsource-id\t[required] The ID of the source resource (volume, snapshot, or image) from which to create the instance.",
				"\t\tboot-index\t[optional] The boot index of the device. Default is 0.",
				"\t\tdelete-on-termination\t[optional] Whether or not to delete the attached volume when the server is delete. Default is false. Options: true, false.",
				"\t\tdestination-type\t[optional] The type that gets created. Options: volume, local.",
				"\t\tvolume-size\t[optional] The size of the volume to create (in gigabytes).",
				"\tExamle: --block-device source-type=image,source-id=bb02b1a3-bc77-4d17-ab5b-421d89850fca,volume-size=100,destination-type=volume,delete-on-termination=false",
			}, "\n"),
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

	if c.IsSet("block-device") {
		bfvMap, err := command.Ctx.CheckKVFlag("block-device")
		if err != nil {
			return err
		}

		sourceID, ok := bfvMap["source-id"]
		if !ok {
			return fmt.Errorf("The source-id key is required when using the --block-device flag.\n")
		}

		sourceTypeRaw, ok := bfvMap["source-type"]
		if !ok {
			return fmt.Errorf("The source-type key is required when using the --block-device flag.\n")
		}
		var sourceType osBFV.SourceType
		switch sourceTypeRaw {
		case "volume", "image", "snapshot":
			sourceType = osBFV.SourceType(sourceTypeRaw)
		default:
			return fmt.Errorf("Invalid value for source-type: %s. Options are: volume, image, snapshot.\n", sourceType)
		}

		bd := osBFV.BlockDevice{
			SourceType: sourceType,
			UUID:       sourceID,
		}

		if volumeSizeRaw, ok := bfvMap["volume-size"]; ok {
			volumeSize, err := strconv.ParseInt(volumeSizeRaw, 10, 16)
			if err != nil {
				return fmt.Errorf("Invalid value for volume-size: %d. Value must be an integer.\n", volumeSize)
			}
			bd.VolumeSize = int(volumeSize)
		}

		if deleteOnTerminationRaw, ok := bfvMap["delete-on-termination"]; ok {
			deleteOnTermination, err := strconv.ParseBool(deleteOnTerminationRaw)
			if err != nil {
				return fmt.Errorf("Invalid value for delete-on-termination: %v. Options are: true, false.\n", deleteOnTermination)
			}
			bd.DeleteOnTermination = deleteOnTermination
		}

		if bootIndexRaw, ok := bfvMap["boot-index"]; ok {
			bootIndex, err := strconv.ParseInt(bootIndexRaw, 10, 8)
			if err != nil {
				return fmt.Errorf("Invalid value for boot-index: %d. Value must be an integer.\n", bootIndex)
			}
			bd.BootIndex = int(bootIndex)
		}

		if destinationType, ok := bfvMap["destination-type"]; ok {
			if destinationType != "volume" && destinationType != "local" {
				return fmt.Errorf("Invalid value for destination-type: %s. Options are: volume, local.\n", destinationType)
			}
			bd.DestinationType = destinationType
		}

		opts.BlockDevice = []osBFV.BlockDevice{bd}
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

	var server *osServers.Server
	var err error
	if len(opts.BlockDevice) > 0 {
		server, err = bfv.Create(command.Ctx.ServiceClient, opts).Extract()
	} else {
		server, err = servers.Create(command.Ctx.ServiceClient, opts).Extract()
	}

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
