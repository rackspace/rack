package handler

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/jrperritt/rack/auth"
)

// Resource is a general resource from Rackspace. This object stores information
// about a single request and response from Rackspace.
type Resource struct {
	// Keys are the fields available to output. These may be limited by the `fields`
	// flag.
	Keys []string
	// Print is a function for custom printing.
	Print *func(*Resource)
	// Params will be the command-specific parameters, such as an instance ID or
	// list options.
	Params interface{}
	// Result will store the result of a single command.
	Result interface{}
	// ErrExit1 will be true if an error was encountered for which the program should
	// exit.
	ErrExit1 bool
	// Err will store any error encountered while processing the command.
	Err error
}

// StreamPipeHandler is an interface that commands implement if they can stream input
// from STDIN.
type StreamPipeHandler interface {
	// PipeHandler is an interface that commands implement if they can accept input
	// from STDIN.
	PipeHandler
	// HandlePipe is a method that commands implement for processing piped input.
	HandleStreamPipe(*Resource) error
}

// PipeHandler is an interface that commands implement if they can accept input
// from STDIN.
type PipeHandler interface {
	// Commander is an interface that all commands will implement.
	Commander
	// HandlePipe is a method that commands implement for processing piped input.
	HandlePipe(*Resource, string) error
	// StdinField is the field that the command accepts on STDIN.
	StdinField() string
}

// Commander is an interface that all commands implement.
type Commander interface {
	// See `Context`.
	Context() *Context
	// Keys returns the keys available for the command output.
	Keys() []string
	// ServiceClientType returns the type of the service client to use.
	ServiceClientType() string
	// HandleFlags processes flags for the command that are relevant for both piped
	// and non-piped commands.
	HandleFlags(*Resource) error
	// HandleSingle contains logic for processing a single resource. This method
	// will be used if input isn't sent to STDIN, so it will contain, for example,
	// logic for handling flags that would be mandatory if otherwise not piped in.
	HandleSingle(*Resource) error
	// Execute executes the command's HTTP request.
	Execute(*Resource)
}

// Handle is the function that handles all commands. It accepts a Commander as
// a parameter, which all commands implement.
func Handle(command Commander) {
	ctx := command.Context()
	ctx.ServiceClientType = command.ServiceClientType()
	ctx.WaitGroup = &sync.WaitGroup{}

	ctx.ListenAndReceive()

	resource := &Resource{
		Keys: command.Keys(),
	}

	err := ctx.CheckArgNum(0)
	if err != nil {
		resource.Err = err
		ctx.ErrExit1(resource)
	}

	client, err := auth.NewClient(ctx.CLIContext, ctx.ServiceClientType)
	if err != nil {
		resource.Err = err
		ctx.ErrExit1(resource)
	}
	ctx.ServiceClient = client

	err = command.HandleFlags(resource)
	if err != nil {
		resource.Err = err
		ctx.ErrExit1(resource)
	}

	// should we expect something on STDIN?
	if ctx.CLIContext.IsSet("stdin") {
		stdinField := ctx.CLIContext.String("stdin")
		// can the command accept input on STDIN?
		if pipeableCommand, ok := command.(PipeHandler); ok {
			// if so, does the given field accept pipeable input?
			if stdinField == pipeableCommand.StdinField() {
				// if so, does the given command and field accept streaming input?
				if streamPipeableCommand, ok := pipeableCommand.(StreamPipeHandler); ok {
					ctx.WaitGroup.Add(1)
					go func() {
						err := streamPipeableCommand.HandleStreamPipe(resource)
						if err != nil {
							resource.Err = fmt.Errorf("Error handling streamable, pipeable command: %s\n", err)
							ctx.Results <- resource
						} else {
							streamPipeableCommand.Execute(resource)
							ctx.Results <- resource
						}
					}()
				} else {
					scanner := bufio.NewScanner(os.Stdin)
					for scanner.Scan() {
						item := scanner.Text()
						ctx.WaitGroup.Add(1)
						go func() {
							err := pipeableCommand.HandlePipe(resource, item)
							if err != nil {
								resource.Err = fmt.Errorf("Error handling pipeable command on %s: %s\n", item, err)
								ctx.Results <- resource
							} else {
								pipeableCommand.Execute(resource)
								ctx.Results <- resource
							}
						}()
					}
					if scanner.Err() != nil {
						resource.Err = scanner.Err()
						ctx.ErrExit1(resource)
					}
				}
			} else {
				resource.Err = fmt.Errorf("Unknown STDIN field: %s\n", stdinField)
				ctx.ErrExit1(resource)
			}
		}
	} else {
		ctx.WaitGroup.Add(1)
		go func() {
			err = command.HandleSingle(resource)
			if err != nil {
				resource.Err = err
				ctx.ErrExit1(resource)
			} else {
				command.Execute(resource)
				ctx.Results <- resource
			}
		}()
	}
	ctx.WaitGroup.Wait()
	ctx.StoreCredentials()
}
