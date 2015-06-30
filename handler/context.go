package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud"
)

// Command is the type that commands have.
type Command struct {
	Ctx *Context
}

// Context is a global context that `rack` uses.
type Context struct {
	// CLIContext is the context that the `cli` library uses. `rack` uses it to
	// access flags.
	CLIContext *cli.Context
	// ServiceClient is the Rackspace service client used to authenticate the user
	// and carry out the requests while processing the command.
	ServiceClient *gophercloud.ServiceClient
	// ServiceClientType is the type of Rackspace service client used (e.g. compute).
	ServiceClientType string
	// WaitGroup is used for synchronizing output.
	WaitGroup *sync.WaitGroup
	// Results is a channel into which commands send results. It allows for streaming
	// output.
	Results chan *Resource
}

// ListenAndReceive creates the Results channel and processes the results that
// come through it before sending them on to `Print`. It is run in a separate
// goroutine from `main`.
func (ctx *Context) ListenAndReceive() {
	ctx.Results = make(chan *Resource)
	go func() {
		for {
			select {
			case resource, ok := <-ctx.Results:
				if !ok {
					ctx.Results = nil
					continue
				}
				if resource.Err != nil {
					ctx.CLIContext.App.Writer = os.Stderr
					resource.Keys = []string{"error"}
					var errorBody string
					switch resource.Err.(type) {
					case *gophercloud.UnexpectedResponseCodeError:
						errBodyRaw := resource.Err.(*gophercloud.UnexpectedResponseCodeError).Body
						errMap := make(map[string]map[string]interface{})
						err := json.Unmarshal(errBodyRaw, &errMap)
						if err != nil {
							errorBody = string(errBodyRaw)
							break
						}
						for _, v := range errMap {
							errorBody = v["message"].(string)
							break
						}
					default:
						errorBody = resource.Err.Error()
					}

					resource.Result = map[string]interface{}{"error": errorBody}
				}
				if resource.Result == nil {
					if args := ctx.CLIContext.Parent().Parent().Args(); len(args) > 0 {
						resource.Result = fmt.Sprintf("Nothing to show. Maybe you'd like to set up some %ss?\n",
							strings.Replace(args[0], "-", " ", -1))
					} else {
						resource.Result = fmt.Sprintf("Nothing to show.\n")
					}
				}
				ctx.Print(resource)
				if resource.ErrExit1 {
					os.Exit(1)
				}
			}
		}
	}()
}

// Print returns the output to the user
func (ctx *Context) Print(resource *Resource) {
	defer ctx.WaitGroup.Done()

	// limit the returned fields if any were given in the `fields` flag
	keys := ctx.limitFields(resource)
	w := ctx.CLIContext.App.Writer
	if ctx.CLIContext.GlobalIsSet("json") || ctx.CLIContext.IsSet("json") {
		switch resource.Result.(type) {
		case map[string]interface{}:
			m := resource.Result.(map[string]interface{})
			output.MetadataJSON(w, m, keys)
		case []map[string]interface{}:
			m := resource.Result.([]map[string]interface{})
			output.ListJSON(w, m, keys)
		case io.Reader:
			if _, ok := resource.Result.(io.ReadCloser); ok {
				defer resource.Result.(io.ReadCloser).Close()
			}
			_, err := io.Copy(w, resource.Result.(io.Reader))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error copying (io.Reader) result: %s\n", err)
			}
		default:
			output.DefaultJSON(w, resource.Result)
		}
	} else if ctx.CLIContext.GlobalIsSet("csv") || ctx.CLIContext.IsSet("csv") {
		switch resource.Result.(type) {
		case map[string]interface{}:
			m := resource.Result.(map[string]interface{})
			output.MetadataCSV(w, m, keys)
		case []map[string]interface{}:
			m := resource.Result.([]map[string]interface{})
			output.ListCSV(w, m, keys)
		case io.Reader:
			if _, ok := resource.Result.(io.ReadCloser); ok {
				defer resource.Result.(io.ReadCloser).Close()
			}
			_, err := io.Copy(w, resource.Result.(io.Reader))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error copying (io.Reader) result: %s\n", err)
			}
		default:
			fmt.Fprintf(w, "%v", resource.Result)
		}
	} else {
		switch resource.Result.(type) {
		case map[string]interface{}:
			m := resource.Result.(map[string]interface{})
			output.MetadataTable(w, m, keys)
		case []map[string]interface{}:
			m := resource.Result.([]map[string]interface{})
			output.ListTable(w, m, keys)
		case io.Reader:
			if _, ok := resource.Result.(io.ReadCloser); ok {
				defer resource.Result.(io.ReadCloser).Close()
			}
			_, err := io.Copy(w, resource.Result.(io.Reader))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error copying (io.Reader) result: %s\n", err)
			}
		default:
			fmt.Fprintf(w, "%v", resource.Result)
		}
	}
}

// limitFields returns only the fields the user specified in the `fields` flag. If
// the flag wasn't provided, all fields are returned.
func (ctx *Context) limitFields(resource *Resource) []string {
	if ctx.CLIContext.IsSet("fields") {
		fields := strings.Split(strings.ToLower(ctx.CLIContext.String("fields")), ",")
		newKeys := []string{}
		for _, key := range resource.Keys {
			if util.Contains(fields, strings.Join(strings.Split(strings.ToLower(key), " "), "-")) {
				newKeys = append(newKeys, key)
			}
		}
		return newKeys
	}
	return resource.Keys
}

// StoreCredentials caches the users auth credentials if available and the `no-cache`
// flag was not provided.
func (ctx *Context) StoreCredentials() {
	// if serviceClient is nil, the HTTP request for the command didn't get sent.
	// don't set cache if the `no-cache` flag is provided
	if ctx.ServiceClient != nil && !ctx.CLIContext.GlobalIsSet("no-cache") && !ctx.CLIContext.IsSet("no-cache") {
		newCacheValue := &auth.CacheItem{
			TokenID:         ctx.ServiceClient.TokenID,
			ServiceEndpoint: ctx.ServiceClient.Endpoint,
		}
		// get auth credentials
		ao, region, err := auth.Credentials(ctx.CLIContext)
		if err == nil {
			// form the cache key
			cacheKey := auth.CacheKey(*ao, region, ctx.ServiceClientType)
			// initialize the cache
			cache := &auth.Cache{}
			// set the cache value to the current values
			_ = cache.SetValue(cacheKey, newCacheValue)
		}
	}
}

// ErrExit1 tells `rack` to print the error and exit.
func (ctx *Context) ErrExit1(resource *Resource) {
	resource.ErrExit1 = true
	ctx.WaitGroup.Add(1)
	ctx.Results <- resource
	ctx.WaitGroup.Wait()
}

// IDOrName is a function for retrieving a resources unique identifier based on
// whether he or she passed an `id` or a `name` flag.
func (ctx *Context) IDOrName(idFromNameFunc func(*gophercloud.ServiceClient, string) (string, error)) (string, error) {
	if ctx.CLIContext.IsSet("id") {
		if ctx.CLIContext.IsSet("name") {
			return "", fmt.Errorf("Only one of either --id or --name may be provided.")
		}
		return ctx.CLIContext.String("id"), nil
	} else if ctx.CLIContext.IsSet("name") {
		name := ctx.CLIContext.String("name")
		id, err := idFromNameFunc(ctx.ServiceClient, name)
		if err != nil {
			return "", fmt.Errorf("Error converting name [%s] to ID: %s", name, err)
		}
		return id, nil
	} else {
		return "", output.ErrMissingFlag{Msg: "One of either --id or --name must be provided."}
	}
}

// CheckArgNum checks that the provided number of arguments has the same
// cardinality as the expected number of arguments.
func (ctx *Context) CheckArgNum(expected int) error {
	argsLen := len(ctx.CLIContext.Args())
	if argsLen != expected {
		return fmt.Errorf("Expected %d args but got %d\nUsage: %s", expected, argsLen, ctx.CLIContext.Command.Usage)
	}
	return nil
}

// CheckFlagsSet checks that the given flag names are set for the command.
func (ctx *Context) CheckFlagsSet(flagNames []string) error {
	for _, flagName := range flagNames {
		if !ctx.CLIContext.IsSet(flagName) {
			return output.ErrMissingFlag{Msg: fmt.Sprintf("--%s is required.", flagName)}
		}
	}
	return nil
}

// CheckKVFlag is a function used for verifying the format of a key-value flag.
func (ctx *Context) CheckKVFlag(flagName string) (map[string]string, error) {
	kv := make(map[string]string)
	kvStrings := strings.Split(ctx.CLIContext.String(flagName), ",")
	for _, kvString := range kvStrings {
		temp := strings.Split(kvString, "=")
		if len(temp) != 2 {
			return nil, output.ErrFlagFormatting{Msg: fmt.Sprintf("Expected key1=value1,key2=value2 format but got %s for --%s.\n", kvString, flagName)}
		}
		kv[temp[0]] = temp[1]
	}
	return kv, nil
}

// CheckStructFlag is a function used for verifying the format of a struct flag.
func (ctx *Context) CheckStructFlag(flagValues []string) ([]map[string]interface{}, error) {
	valSliceMap := make([]map[string]interface{}, len(flagValues))
	for i, flagValue := range flagValues {
		kvStrings := strings.Split(flagValue, ",")
		m := make(map[string]interface{})
		for _, kvString := range kvStrings {
			temp := strings.Split(kvString, "=")
			if len(temp) != 2 {
				return nil, output.ErrFlagFormatting{Msg: fmt.Sprintf("Expected key1=value1,key2=value2 format but got %s.\n", kvString)}
			}
			m[temp[0]] = temp[1]
		}
		valSliceMap[i] = m
	}
	return valSliceMap, nil
}
