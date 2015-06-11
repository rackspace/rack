package output

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// Print prints the results of the CLI command. This function is designed to centralize
// printing command output and to accomodate the way the Rackspace API is designed.
// When performing actions on resources that provide a response body, there are 2
// basic formats the response body can come in: a slice of items (which would be
// returned from a `list` operation) or a map of items (which would be returned
// from operations like `get` and `create`). In addition, to these 2 formats,
// sometimes it is more convenient for the user to just get back the most important
// piece of information from a command. For example, a `rack servers keypair get`
// command will return just the raw public key. Even though there are other fields
// returned in the response body, the public key is almost always what the user
// wants from the command and printing out just the public key makes it easy to
// copy (and paste).
//
// This function accepts 3 parameters:
// c) a pointer to the cli.Context: this is for checking the flags to determine the
// 		output format.
// f) a pointer to a function that returns an `interface{}`: this (pointer to a)
// 		function returns the data from the response body that gets returned from the
//		command request. It returns an `interface{}` to accomodate the different data
//		types that can be printed (e.g. `map[string]interface{}`, `[]map[string]interface{}`,
//		`string`, etc). Some commands can print out the data that gets returned as-is,
//		for example if all the values are `string` or `int`. In these cases, a simple
//		operation that converts the value(s) from the response body into `map[string]interface{}`
//		(or `[]map[string]interface{}`) is sufficient. In such cases, the parameter `f`
//		in the `Print` function might look like this:
//
//		f := func() interface{} {
//			// o would be defined above and is the raw (or type casted) value
//			// returned in the response body.
//			// This function returns an `interface{}` that has a `map[string]interface{}`
//			// type.
//			return structs.Map(o)
//		}
//
//		However, in some cases it is necessary to massage the data to get it into
//		a presentable form, like if the response body contains nested maps and one
//		of the nested maps has a field that you want to output. Below is an example
//		of a function for a `list` command:
//
//		f := func() interface{} {
//			// o would be defined above and is the raw (or type casted) value
//			// returned in the response body
//			m := make([]map[string]interface{}, len(o))
//			for j, kp := range o {
//				m[j] = structs.Map(kp)
//			}
//			// This function returns an `interface{}` that has a `[]map[string]interface{}`
//			// type.
//			return m
//		}
//
//		Still yet, there may be cases when you just want a single piece of raw data
//		(like in the case of `rack servers keypair get`). In these cases, you can
//		choose to return a single field from the response body as a string:
//
//			f := func() interface{} {
//				// o would be defined above and is the raw (or type casted) value
//				// returned in the response body
//				m := structs.Map(o)
//				// This function returns an `interface{}` that has a `string` type.
//				return m["PublicKey"]
//			}
//
//		Regardless of how the function looks, it is created within the command
// 		function as a closure around the data from response body.
//
// keys) a slice of strings: this slice contains the header values to print out for
// 		the tabular and csv formats.
func Print(c *cli.Context, f *func() interface{}, keys []string) {
	i := (*f)()
	if c.IsSet("json") {
		jsonOut(i)
		return
	}
	if c.IsSet("csv") {
		switch i.(type) {
		case map[string]interface{}:
			m := i.(map[string]interface{})
			metadataCSV(c, m, keys)
		case []map[string]interface{}:
			m := i.([]map[string]interface{})
			listCSV(c, m, keys)
		default:
			fmt.Fprintf(c.App.Writer, "%v", i)
		}
		return
	}
	switch i.(type) {
	case map[string]interface{}:
		m := i.(map[string]interface{})
		metadataTable(c, m, keys)
	case []map[string]interface{}:
		m := i.([]map[string]interface{})
		listTable(c, m, keys)
	default:
		fmt.Fprintf(c.App.Writer, "%v", i)
	}
}
