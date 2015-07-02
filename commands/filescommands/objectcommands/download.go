package objectcommands

import (
	"io/ioutil"

	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
	"github.com/jrperritt/rack/util"
)

var download = cli.Command{
	Name:        "download",
	Usage:       util.Usage(commandPrefix, "download", "--container <containerName> --name <objectName>"),
	Description: "Downloads an object",
	Action:      actionDownload,
	Flags:       util.CommandFlags(flagsDownload, keysDownload),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDownload, keysDownload))
	},
}

func flagsDownload() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container containing the object to download",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the object to download",
		},
	}
}

var keysDownload = []string{}

type paramsDownload struct {
	container string
	object    string
}

type commandDownload handler.Command

func actionDownload(c *cli.Context) {
	command := &commandDownload{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandDownload) Context() *handler.Context {
	return command.Ctx
}

func (command *commandDownload) Keys() []string {
	return keysDownload
}

func (command *commandDownload) ServiceClientType() string {
	return serviceClientType
}

func (command *commandDownload) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"container", "name"})
	if err != nil {
		return err
	}
	container := command.Ctx.CLIContext.String("container")
	object := command.Ctx.CLIContext.String("name")
	resource.Params = &paramsDownload{
		container: container,
		object:    object,
	}
	return nil
}

func (command *commandDownload) Execute(resource *handler.Resource) {
	containerName := resource.Params.(*paramsDownload).container
	objectName := resource.Params.(*paramsDownload).object
	rawResponse := objects.Download(command.Ctx.ServiceClient, containerName, objectName, nil)
	if rawResponse.Err != nil {
		resource.Err = rawResponse.Err
		return
	}
	if command.Ctx.OutputFormat == "json" {
		bytes, err := ioutil.ReadAll(rawResponse.Body)
		if err != nil {
			resource.Err = err
			return
		}
		resource.Result = string(bytes)
	} else {
		resource.Result = rawResponse.Body
	}
}
