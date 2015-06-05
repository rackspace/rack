package util

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

// PrintError is used for printing information when an error is encountered.
func PrintError(c *cli.Context, err error) {
	w := c.App.Writer
	switch err.(type) {
	case ErrMissingFlag, ErrFlagFormatting:
		fmt.Fprintf(w, "%s", err.Error())
		fmt.Fprintf(w, "Usage: %s\n", c.Command.Usage)
	}
	os.Exit(1)
}

// ErrMissingFlagPrefix is the prefix for when a required flag is missing.
var ErrMissingFlagPrefix = "Missing flag:"

// ErrMissingFlag is used when a user doesn't provide a required flag.
type ErrMissingFlag struct {
	Msg string
}

func (e ErrMissingFlag) Error() string {
	return fmt.Sprintf("%s %s\n", ErrMissingFlagPrefix, e.Msg)
}

// ErrFlagFormatting is the prefix for when a flag's format is invalid.
var ErrFlagFormattingPrefix = "Invalid flag formatting:"

// ErrFlagFormatting is used when a flag's format is invalid.
type ErrFlagFormatting struct {
	Msg string
}

func (e ErrFlagFormatting) Error() string {
	return fmt.Sprintf("%s %s\n", ErrFlagFormattingPrefix, e.Msg)
}
