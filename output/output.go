package output

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
)

// JSON prints results in JSON format.
func JSON(i interface{}) {
	j, _ := json.Marshal(i)
	fmt.Println(string(j))
}

// Print prints the results of the CLI command.
func Print(c *cli.Context, i interface{}, table func(*cli.Context, interface{})) {
	if c.IsSet("json") {
		JSON(i)
		return
	}
	table(c, i)
	return
}
