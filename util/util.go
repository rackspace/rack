package util

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/codegangsta/cli"
)

// Name is the name of the CLI
var Name = "rack"

// Version is the current CLI version
var Version = "0.0.0-dev"

// UserAgent is the user-agent used for each HTTP request
var UserAgent = fmt.Sprintf("%s-%s/%s", "rackcli", runtime.GOOS, Version)

// Usage return a string that specifies how to call a particular command.
func Usage(commandPrefix, action, mandatoryFlags string) string {
	return fmt.Sprintf("%s [GLOBALS] %s %s %s [OPTIONS]", Name, commandPrefix, action, mandatoryFlags)
}

// Contains checks whether a given string is in a provided slice of strings.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// RackDir returns the location of the `rack` directory. This directory is for
// storing `rack`-specific information such as the cache or a config file.
func RackDir() (string, error) {
	var homeDir string
	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH") // Windows
		if homeDir == "" {
			homeDir = os.Getenv("USERPROFILE") // Windows
		}
	} else {
		homeDir = os.Getenv("HOME") // *nix
	}
	if homeDir == "" {
		return "", errors.New("User home directory not found.")
	}
	dirpath := path.Join(homeDir, ".rack")
	err := os.MkdirAll(dirpath, 0744)
	return dirpath, err
}

// IDAndNameFlags are flags for commands that allow either an ID or a name.
var IDAndNameFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional] The ID of the resource",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional] The name of the resource",
	},
}

// IDOrNameUsage returns flag usage information for resources that allow either
// an ID or a name.
func IDOrNameUsage(resource string) string {
	return fmt.Sprintf("[--id <%sID> | --name <%sName>]", resource, resource)
}
