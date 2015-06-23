package output

import "fmt"

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

// ErrArgsFlagPrefix is the prefix for when a flag's argument is invalid.
var ErrArgsFlagPrefix = "Argument error:"

// ErrArgs is used when a flag's arguments are invalid.
type ErrArgs struct {
	Msg string
}

func (e ErrArgs) Error() string {
	return fmt.Sprintf("%s %s\n", ErrArgsFlagPrefix, e.Msg)
}
