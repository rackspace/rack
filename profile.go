package main

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/gopkg.in/ini.v1"
)

func profileCommandsGet() []cli.Command {
	return []cli.Command{
		{
			Name:        "activate",
			Description: "Activate a profile",
			Usage:       "rack profile activate --name <profile-name>",
			Action:      profileActivate,
			Flags:       profileFlagsActivate,
			BashComplete: func(c *cli.Context) {
				commandoptions.CompleteFlags(profileFlagsActivate)
			},
		},
		{
			Name:        "list",
			Description: "List profile information",
			Usage:       "rack profile list",
			Action:      profileList,
		},
	}
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

	configFileLoc, err := configFileLocation()
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
	configFileLoc, err := configFileLocation()
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
