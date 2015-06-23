package objectcommands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osObjects "github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
)

var upload = cli.Command{
	Name:        "upload",
	Usage:       util.Usage(commandPrefix, "upload", "--container <containerName> --name <objectName>"),
	Description: "Uploads an object",
	Action:      actionUpload,
	Flags:       util.CommandFlags(flagsUpload, keysUpload),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsUpload, keysUpload))
	},
}

func flagsUpload() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container to upload the object upload",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name the object should have in the Cloud Files container",
		},
		cli.StringFlag{
			Name:  "content",
			Usage: "[optional; required if `file` or `stdin` isn't provided] The string contents to upload",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "[optional; required if `content` or `stdin` isn't provided] The file name containing the contents to upload",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `file` or `content` isn't provided] The field being piped to STDIN, if any. Valid values are: file",
		},
	}
}

var keysUpload = []string{}

type paramsUpload struct {
	container string
	object    string
	content   string
	fileName  string
	opts      osObjects.CreateOpts
}

type commandUpload handler.Command

func actionUpload(c *cli.Context) {
	command := &commandUpload{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUpload) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUpload) Keys() []string {
	return keysUpload
}

func (command *commandUpload) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUpload) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"container", "name"})
	if err != nil {
		return err
	}

	container := command.Ctx.CLIContext.String("container")
	object := command.Ctx.CLIContext.String("name")
	resource.Params = &paramsUpload{
		container: container,
		object:    object,
	}
	return nil
}

func (command *commandUpload) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsUpload).fileName = item
	return nil
}

func (command *commandUpload) HandleSingle(resource *handler.Resource) error {
	if command.Ctx.CLIContext.IsSet("file") {
		resource.Params.(*paramsUpload).fileName = command.Ctx.CLIContext.String("file")
	} else if command.Ctx.CLIContext.IsSet("content") {
		resource.Params.(*paramsUpload).content = command.Ctx.CLIContext.String("content")
	} else {
		return fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN.")
	}
	return nil
}

func (command *commandUpload) Execute(resource *handler.Resource) {
	containerName := resource.Params.(*paramsUpload).container
	objectName := resource.Params.(*paramsUpload).object
	opts := resource.Params.(*paramsUpload).opts
	var readSeeker io.ReadSeeker
	var err error
	if fileName := resource.Params.(*paramsUpload).fileName; fileName != "" {
		// this file will be closed by Gophercloud, if not closed before then
		readSeeker, err = os.Open(fileName)
		if err != nil {
			resource.Err = err
			return
		}
	} else {
		content := resource.Params.(*paramsUpload).content
		readSeeker = strings.NewReader(content)
	}
	rawResponse := objects.Create(command.Ctx.ServiceClient, containerName, objectName, readSeeker, opts)
	if rawResponse.Err != nil {
		resource.Err = rawResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully uploaded object [%s] to container [%s]\n", objectName, containerName)
}

func (command *commandUpload) StdinField() string {
	return "file"
}
