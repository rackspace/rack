package commandoptions

import (
	"fmt"
	"path"

	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/gopkg.in/ini.v1"
	"github.com/jrperritt/rack/util"
)

func ConfigFile(c *cli.Context, have map[string]Cred, need map[string]string) error {
	var profile string
	if c.GlobalIsSet("profile") {
		profile = c.GlobalString("profile")
	} else if c.IsSet("profile") {
		profile = c.String("profile")
	}

	section, err := ProfileSection(profile)
	if err != nil {
		return err
	}

	if section == nil {
		return nil
	}

	for opt := range need {
		if val := section.Key(opt).String(); val != "" {
			have[opt] = Cred{Value: val, From: fmt.Sprintf("config file (profile: %s)", section.Name())}
			delete(need, opt)
		}
	}

	if profile != "" {
		section, err := ProfileSection("")
		if err != nil {
			return err
		}

		for opt := range need {
			if val := section.Key(opt).String(); val != "" {
				have[opt] = Cred{Value: val, From: fmt.Sprintf("config file (profile: default)")}
				delete(need, opt)
			}
		}
	}

	return nil
}

func ProfileSection(profile string) (*ini.Section, error) {
	dir, err := util.RackDir()
	if err != nil {
		// return fmt.Errorf("Error retrieving rack directory: %s\n", err)
		return nil, nil
	}
	f := path.Join(dir, "config")
	cfg, err := ini.Load(f)
	if err != nil {
		// return fmt.Errorf("Error loading config file: %s\n", err)
		return nil, nil
	}
	cfg.BlockMode = false
	section, err := cfg.GetSection(profile)
	if err != nil && profile != "" {
		return nil, fmt.Errorf("Invalid config file profile: %s\n", profile)
	}
	return section, nil
}
