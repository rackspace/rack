package util

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"

	"gopkg.in/ini.v1"
)

// Name is the name of the CLI
var Name = "rack"

// UserAgent is the user-agent used for each HTTP request
var UserAgent = fmt.Sprintf("%s-%s/%s", "rackcli", runtime.GOOS, Version)

// Usage return a string that specifies how to call a particular command.
func Usage(commandPrefix, action, mandatoryFlags string) string {
	return fmt.Sprintf("%s %s %s %s [flags]", Name, commandPrefix, action, mandatoryFlags)
}

// RemoveFromList removes an element from a slice and returns the slice.
func RemoveFromList(list []string, item string) []string {
	for i, element := range list {
		if element == item {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	return list
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
	homeDir, err := HomeDir()
	if err != nil {
		return "", err
	}
	dirpath := path.Join(homeDir, ".rack")
	err = os.MkdirAll(dirpath, 0744)
	return dirpath, err
}

// HomeDir returns the user's home directory, which is platform-dependent.
func HomeDir() (string, error) {
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
	return homeDir, nil
}

func ConfigFileLocation() (string, error) {
	dir, err := RackDir()
	if err != nil {
		return "", fmt.Errorf("Error fetching config directory: %s", err)
	}
	filepath := path.Join(dir, "config")
	// check if the config file exists
	if _, err := os.Stat(filepath); err == nil {
		return filepath, nil
	}
	// create the config file if it doesn't already exist
	f, err := os.Create(filepath)
	defer f.Close()
	return filepath, err
}

func CanActivateProfile() bool {
	configFileLoc, err := ConfigFileLocation()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error to determining config file location: %s\n", err)
		return false
	}

	cfg, err := ini.Load(configFileLoc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %s\n", err)
		return false
	}

	chosenSection, err := cfg.GetSection("DEFAULT")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Section [%s] doesn't exist in config file\n", "DEFAULT")
		return false
	}

	if admin, ok := chosenSection.KeysHash()["enable-profile-activate"]; ok && admin == "true" {
		return true
	}

	return false
}

// Pluralize will plurarize a given noun according to its number. For example,
// 0 servers were deleted; 1 account updated.
func Pluralize(noun string, count int64) string {
	if count != 1 {
		noun += "s"
	}
	return noun
}

// The following are both set during build time

// Version is the current CLI version
var Version string

// Commit is the commit this build comes from
var Commit string
