package auth

import (
	"fmt"
	"path"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/util"
	"gopkg.in/ini.v1"
)

func configfile(c *cli.Context, have map[string]string, need map[string]string) {
	dir, err := util.RackDir()
	if err != nil {
		fmt.Fprint(c.App.Writer, err)
		return
	}
	f := path.Join(dir, "config")
	cfg, err := ini.Load(f)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}
	cfg.BlockMode = false
	var profile string
	if c.GlobalIsSet("profile") {
		profile = c.GlobalString("profile")
	}
	section, err := cfg.GetSection(profile)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	for opt := range need {
		if val := section.Key(opt).String(); val != "" {
			have[opt] = val
			delete(need, opt)
		}
	}
}
