package util

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
)

// Name is the name of the CLI
var Name = "rack"

// UserAgent is the user-agent used for each HTTP request
var UserAgent = fmt.Sprintf("%s-%s/%s", "rackcli", runtime.GOOS, Version)

// Usage return a string that specifies how to call a particular command.
func Usage(commandPrefix, action, mandatoryFlags string) string {
	return fmt.Sprintf("%s %s %s %s [OPTIONS]", Name, commandPrefix, action, mandatoryFlags)
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

// Pluralize will plurarize a given noun according to its number. For example,
// 0 servers were deleted; 1 account updated.
func Pluralize(noun string, count int64) string {
	if count != 1 {
		noun += "s"
	}
	return noun
}

// Version is the current CLI version
var Version = "0.0.0-dev"
var Commit = "07e635a690a073eed5969581794738170c6be90c"
