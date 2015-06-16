package util

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// Error is used for printing information when an error is encountered.
func Error(c *cli.Context, e error) error {
	switch e.(type) {
	case ErrMissingFlag, ErrFlagFormatting, ErrArgs:
		return fmt.Errorf("%s\nUsage: %s\n", e.Error(), c.Command.Usage)
	}
	return fmt.Errorf("%s\n", e)
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

var ErrArgsFlagPrefix = "Argument error:"

type ErrArgs struct {
	Msg string
}

func (e ErrArgs) Error() string {
	return fmt.Sprintf("%s %s\n", ErrArgsFlagPrefix, e.Msg)
}
