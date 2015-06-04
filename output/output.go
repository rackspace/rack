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
	switch c.String("format") {
	case "json":
		JSON(i)
		return
	default:
		table(c, i)
		return
	}
}
