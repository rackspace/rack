package containercommands

import (
	"fmt"
	"sync"
	"time"

	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	"github.com/jrperritt/rack/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--name <containerName> | --stdin name]"),
	Description: "Deletes a container",
	Action:      actionDelete,
	Flags:       util.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container",
		},
		cli.BoolFlag{
			Name:  "purge",
			Usage: "[optional] If set, this command will delete all objects in the container, and then delete the container.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	container string
}

type commandDelete handler.Command

func actionDelete(c *cli.Context) {
	command := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandDelete) Context() *handler.Context {
	return command.Ctx
}

func (command *commandDelete) Keys() []string {
	return keysDelete
}

func (command *commandDelete) ServiceClientType() string {
	return serviceClientType
}

func (command *commandDelete) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsDelete{}
	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsDelete).container = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).container = command.Ctx.CLIContext.String("name")
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsDelete)
	containerName := params.container
	if command.Ctx.CLIContext.IsSet("purge") {
		allPages, err := objects.List(command.Ctx.ServiceClient, containerName, nil).AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		objectNames, err := objects.ExtractNames(allPages)
		if err != nil {
			resource.Err = err
			return
		}
		wg := &sync.WaitGroup{}
		for _, objectName := range objectNames {
			wg.Add(1)
			go func(objectName string) {
				defer wg.Done()
				rawResponse := objects.Delete(command.Ctx.ServiceClient, containerName, objectName, nil)
				if rawResponse.Err != nil {
					resource.Err = rawResponse.Err
					return
				}
			}(objectName)
		}
		wg.Wait()
		numTimesChecked := 0
		for {
			allPages, err := objects.List(command.Ctx.ServiceClient, containerName, nil).AllPages()
			if err != nil {
				resource.Err = err
				return
			}
			objectNames, err := objects.ExtractNames(allPages)
			if err != nil {
				resource.Err = err
				return
			}
			if len(objectNames) == 0 {
				break
			}
			numTimesChecked++
			if numTimesChecked == 60 {
				resource.Err = fmt.Errorf("Purging objects from container [%s] timed out. There are still %d object left.\n", containerName, len(objectNames))
			}
			time.Sleep(5 * time.Second)
		}
	}
	rawResponse := containers.Delete(command.Ctx.ServiceClient, containerName)
	if rawResponse.Err != nil {
		resource.Err = rawResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted container [%s]\n", containerName)
}

func (command *commandDelete) StdinField() string {
	return "name"
}
