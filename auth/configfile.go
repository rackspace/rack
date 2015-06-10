package auth

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"gopkg.in/ini.v1"
)

func configfile(c *cli.Context, have map[string]string, need map[string]string) {
	dir, err := rackDir()
	if err != nil {
		fmt.Fprint(c.App.Writer, err)
		return
	}
	f := path.Join(dir, "creds")
	cfg, err := ini.Load(f)
	fmt.Printf("cfg: %+v\n", cfg)
}

func rackDir() (string, error) {
	homeDir := os.Getenv("HOME") // *nix
	if homeDir == "" {           // Windows
		homeDir = os.Getenv("USERPROFILE")
	}
	if homeDir == "" {
		return "", errors.New("User home directory not found.")
	}

	return path.Join(homeDir, ".rack"), nil
}
