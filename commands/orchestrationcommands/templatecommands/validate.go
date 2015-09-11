package templatecommands

import (
	"io/ioutil"
	"path/filepath"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStackTemplates "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacktemplates"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacktemplates"
	"github.com/rackspace/rack/util"
)

var validate = cli.Command{
	Name:        "validate",
	Usage:       util.Usage(commandPrefix, "validate", "[--template <templateFile> | --template-url <templateURL>]"),
	Description: "Validate a specified template",
	Action:      actionValidate,
	Flags:       commandoptions.CommandFlags(flagsValidate, keysValidate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsValidate, keysValidate))
	},
}

func flagsValidate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "template",
			Usage: "[optional; required if `template-url` isn't provided] The path to template file.",
		},
		cli.StringFlag{
			Name:  "template-url",
			Usage: "[optional; required if `template` isn't provided] The url to template.",
		},
	}
}

type paramsValidate struct {
	opts *osStackTemplates.ValidateOpts
}

var keysValidate = []string{"Description", "Parameters", "ParameterGroups"}

type commandValidate handler.Command

func actionValidate(c *cli.Context) {
	command := &commandValidate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandValidate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandValidate) Keys() []string {
	return keysValidate
}

func (command *commandValidate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandValidate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	opts := osStackTemplates.ValidateOpts{}

	// check if either template url or template file is set
	if c.IsSet("template") {
		abs, err := filepath.Abs(c.String("template"))
		if err != nil {
			return err
		}
		template, err := ioutil.ReadFile(abs)
		if err != nil {
			return err
		}
		opts.Template = string(template)
	} else if c.IsSet("template-url") {
		opts.TemplateURL = c.String("template-url")
	}

	resource.Params = &paramsValidate{
		opts: &opts,
	}
	return nil
}

func (command *commandValidate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsValidate).opts

	validateResult, err := stacktemplates.Validate(command.Ctx.ServiceClient, params).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = string(validateResult)
}
