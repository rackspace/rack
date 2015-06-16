package util

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/rackspace/gophercloud"
)

// Name is the name of the CLI
var name = "rack"

// Version is the current CLI version
var Version = "0.0.0-dev"

// Usage return a string that specifies how to call a particular command.
func Usage(commandPrefix, action, mandatoryFlags string) string {
	return fmt.Sprintf("%s [globals] %s %s %s [flags]", name, commandPrefix, action, mandatoryFlags)
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
	homeDir := os.Getenv("HOME") // *nix
	if homeDir == "" {           // Windows
		homeDir = os.Getenv("USERPROFILE")
	}
	if homeDir == "" {
		return "", errors.New("User home directory not found.")
	}

	return path.Join(homeDir, ".rack"), nil
}

// CheckArgNum checks that the provided number of arguments has the same
// cardinality as the expected number of arguments.
func CheckArgNum(c *cli.Context, expected int) error {
	argsLen := len(c.Args())
	if argsLen != expected {
		return fmt.Errorf("Expected %d args but got %d\nUsage: %s", expected, argsLen, c.Command.Usage)
	}
	return nil
}

// CheckFlagsSet checks that the given flag names are set for the command.
func CheckFlagsSet(c *cli.Context, flagNames []string) error {
	for _, flagName := range flagNames {
		if !c.IsSet(flagName) {
			return Error(c, ErrMissingFlag{
				Msg: fmt.Sprintf("--%s is required.", flagName),
			})
		}
	}
	return nil
}

// IDOrName is a function for retrieving a resources unique identifier based on
// whether he or she passed an `id` or a `name` flag.
func IDOrName(c *cli.Context, client *gophercloud.ServiceClient, idFromName func(*gophercloud.ServiceClient, string) (string, error)) (string, error) {
	if c.IsSet("id") {
		return c.String("id"), nil
	} else if c.IsSet("name") {
		name := c.String("name")
		id, err := idFromName(client, name)
		if err != nil {
			return "", fmt.Errorf("Error converting name [%s] to ID: %s", name, err)
		}
		return id, nil
	} else {
		return "", Error(c, ErrMissingFlag{
			Msg: "One of either --id or --name must be provided.",
		})
	}
}

// IDAndNameFlags are flags for commands that allow either an ID or a name.
var IDAndNameFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional; required if 'name' is not provided] The ID of the resource",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if 'id' is not provided] The name of the resource",
	},
}

// IDOrNameUsage returns flag usage information for resources that allow either
// an ID or a name.
func IDOrNameUsage(resource string) string {
	return fmt.Sprintf("[--id <%sID> | --name <%sName>]", resource, resource)
}
