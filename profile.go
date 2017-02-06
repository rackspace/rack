package main

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/codegangsta/cli"
	"gopkg.in/ini.v1"
	"github.com/rackspace/rack/util"
)

var commandActivate = cli.Command{
	Name: "activate",
	Description: "Activate a profile. Activating a profile will have the following\n" +
		"\torder for the way that `rack` looks for command configuration values:\n" +
		"\t1. command-line options\n" +
		"\t2. the active profile\n" +
		"\t3. the default profile\n" +
		"\t4. environment variables\n" +
		"\n" +
		"\tNOTE: The safest way to use `rack` is by always explicitly providing\n" +
		"\tconfiguration values (like `--profile`) as command-line options. Running\n" +
		"\ta command without knowing which profile is active can result in unintended\n" +
		"\tconsequences.",
	Usage:  "rack profile activate --name <profile-name>",
	Action: profileActivate,
	Flags:  profileFlagsActivate,
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(profileFlagsActivate)
	},
}

var commandList = cli.Command{
	Name:        "list",
	Description: "List profile information",
	Usage:       "rack profile list",
	Action:      profileList,
}

var commandsProfile = []cli.Command{
	commandList,
}

var commandsProfileAdmin = []cli.Command{
	commandActivate,
}

func profileCommandsGet(canActivateProfile bool) []cli.Command {
	if canActivateProfile {
		return append(commandsProfile, commandsProfileAdmin...)
	}
	return commandsProfile
}

var profileFlagsActivate = []cli.Flag{
	cli.StringFlag{
		Name:  "name",
		Usage: "[required] The name of the profile to activate",
	},
}

func profileActivate(c *cli.Context) {
	if !c.IsSet("name") {
		fmt.Fprintln(c.App.Writer, "You must provide the profile name with the `name` flag")
		return
	}

	profileName := c.String("name")

	configFileLoc, err := util.ConfigFileLocation()
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Error to determining config file location: %s\n", err)
		return
	}

	cfg, err := ini.Load(configFileLoc)
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Error loading config file: %s\n", err)
		return
	}

	chosenSection, err := cfg.GetSection(profileName)
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Section [%s] doesn't exist in config file\n", profileName)
		return
	}

	sections := cfg.Sections()
	for _, section := range sections {
		section.DeleteKey("enabled")
	}

	chosenSection.Key("enabled").SetValue("true")

	err = cfg.SaveTo(configFileLoc)
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Error saving config file: %s\n", err)
		return
	}

	fmt.Fprintf(c.App.Writer, "Successfully activated profile [%s]\n", profileName)
}

func profileList(c *cli.Context) {
	configFileLoc, err := util.ConfigFileLocation()
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Error to determining config file location: %s\n", err)
		return
	}

	cfg, err := ini.Load(configFileLoc)
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Error loading config file: %s\n", err)
		return
	}

	sections := cfg.Sections()
	for _, section := range sections {
		fmt.Fprintf(c.App.Writer, "[%s]\n", section.Name())
		for k, v := range section.KeysHash() {
			fmt.Fprintf(c.App.Writer, "%s = %s\n", k, v)
		}
		fmt.Fprintln(c.App.Writer)
		fmt.Fprintln(c.App.Writer)
	}
}
