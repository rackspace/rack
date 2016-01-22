package instancecommands

import (
	"errors"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	db "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/databases"
	os "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/instances"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--name <instanceName> --flavor <flavorId> --size <instanceSize>"),
	Description: "Creates a new database instance",
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
			Usage: "[required] The name that the database instance should have.",
		},
		cli.StringFlag{
			Name:  "flavor",
			Usage: "[required] The flavor ID that the database instance will be based on",
		},
		cli.IntFlag{
			Name:  "size",
			Usage: "[required] The disk space that will be allocated for the database instance in GB",
		},
		cli.StringSliceFlag{
			Name: "database",
			Usage: "[optional] A database to be added on the database instance. Should " +
				"be formatted like so: <name>:<charset>:<collate>. Only <name> is required. " +
				"If <charset> is omitted, `utf8` is used. If <collate> is omitted, `utf8_general_ci` is used. " +
				"Some examples would be new_instance:utf8:utf8_general_ci, new_instance:utf8 or new_instance::utf8_general_ci",
		},
		cli.StringSliceFlag{
			Name: "user",
			Usage: "[optional] A user to be added to the database instance. Should be " +
				"formatted like so: <username>:<password>:<databases>:<host> where <databases> is" +
				"a comma-delimeted list of databases like `db1,db2,db3`. Only <username> " +
				"and <password> are required",
		},
		cli.StringFlag{
			Name:  "config-id",
			Usage: "[optional] UUID of the configuration group to associate with the instance.",
		},
		cli.StringFlag{
			Name:  "datastore-type",
			Usage: "[optional] The type of the datastore, for example `mysql`",
		},
		cli.StringFlag{
			Name:  "datastore-version",
			Usage: "[optional] The specific version of a datastore, for example `5.6`",
		},
		cli.StringFlag{
			Name:  "restore-point",
			Usage: "[optional] Specifies the backup ID from which to restore the database instance.",
		},
		cli.StringFlag{
			Name:  "replica-of",
			Usage: "[optional] The UUID of the instance that this new database instance will replicate",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the instance is available.",
		},
	}
}

var keysCreate = []string{"ID", "Hostname"}

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

type paramsCreate struct {
	wait bool
	opts *instances.CreateOpts
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	wait := false
	if c.IsSet("wait-for-completion") {
		wait = true
	}

	opts := &instances.CreateOpts{
		FlavorRef: c.String("flavor"),
		Name:      c.String("name"),
		Size:      c.Int("size"),
	}

	// add user values
	usersSlice := c.StringSlice("user")
	if usersSlice != nil {
		var usrs []users.CreateOpts
		for _, userStr := range usersSlice {
			user := strings.Split(userStr, ":")
			if len(user) < 2 {
				return errors.New("Incorrect usage for --user. Correct usage is: " +
					"<username>:<password>:<db1,db2,db3>:<host>. <username> and <password> values " +
					"are required as a minimum")
			}

			userOpts := users.CreateOpts{
				Name:     user[0],
				Password: user[1],
			}

			if len(user) >= 3 {
				userOpts.Databases = db.BatchCreateOpts{}
				for _, dbName := range strings.Split(user[2], ",") {
					userOpts.Databases = append(userOpts.Databases, db.CreateOpts{
						Name: dbName,
					})
				}
			}

			if len(user) >= 4 {
				userOpts.Host = user[3]
			}

			usrs = append(usrs, userOpts)
		}
		opts.Users = users.BatchCreateOpts(usrs)
	}

	// add db values
	dbsSlice := c.StringSlice("database")
	if usersSlice != nil {
		var dbs []db.CreateOpts
		for _, dbStr := range dbsSlice {
			dbSlice := strings.Split(dbStr, ":")
			dbOpts := db.CreateOpts{Name: dbSlice[0]}
			if len(dbSlice) >= 2 {
				dbOpts.CharSet = dbSlice[1]
			}
			if len(dbSlice) >= 3 {
				dbOpts.Collate = dbSlice[2]
			}
			dbs = append(dbs, dbOpts)
		}
		opts.Databases = db.BatchCreateOpts(dbs)
	}

	if c.IsSet("config-id") {
		opts.ConfigID = c.String("config-id")
	}

	datastoreErr := errors.New("If specifying a datastore, you must provide both a version and a type")

	if c.IsSet("datastore-version") {
		if !c.IsSet("datastore-type") {
			return datastoreErr
		}
	}

	if c.IsSet("datastore-type") {
		if !c.IsSet("datastore-version") {
			return datastoreErr
		}
		opts.Datastore = &os.DatastoreOpts{
			Version: c.String("datastore-version"),
			Type:    c.String("datastore-type"),
		}
	}

	if c.IsSet("restore-point") {
		opts.RestorePoint = c.String("restore-point")
	}

	if c.IsSet("replica-of") {
		opts.ReplicaOf = c.String("replica-of")
	}

	resource.Params = &paramsCreate{
		wait: wait,
		opts: opts,
	}
	return nil
}

func (command *commandCreate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsCreate).opts.Name = item
	return nil
}

func (command *commandCreate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name", "flavor", "size"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsCreate).opts.Name = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts

	instance, err := instances.Create(command.Ctx.ServiceClient, opts).Extract()

	if resource.Params.(*paramsCreate).wait {
		err = gophercloud.WaitFor(600, func() (bool, error) {
			inst, err := instances.Get(command.Ctx.ServiceClient, instance.ID).Extract()
			if err != nil {
				return false, err
			}
			if inst.Status == "ACTIVE" {
				return true, nil
			}
			return false, nil
		})
	}

	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = singleInstance(instance)
}
