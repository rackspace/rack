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

// ErrMissingFlag is used when a user doesn't provide a required flag.
type ErrMissingFlag struct {
	Msg string
}

func (e ErrMissingFlag) Error() string {
	return fmt.Sprintf("Missing flag: %s\n", e.Msg)
}

// ErrFlagFormatting is used when a flag's format is invalid.
type ErrFlagFormatting struct {
	Msg string
}

func (e ErrFlagFormatting) Error() string {
	return fmt.Sprintf("Invalid flag formatting: %s\n", e.Msg)
}
