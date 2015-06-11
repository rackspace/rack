package output

import (
	"encoding/csv"
	"fmt"

	"github.com/codegangsta/cli"
)

func metadataCSV(c *cli.Context, m map[string]interface{}, keys []string) {
	w := csv.NewWriter(c.App.Writer)
	w.Write([]string{"PROPERTY", "VALUE"})
	for _, key := range keys {
		w.Write([]string{key, fmt.Sprint(m[key])})
	}
	w.Flush()
}

func listCSV(c *cli.Context, many []map[string]interface{}, keys []string) {
	w := csv.NewWriter(c.App.Writer)
	w.Write(keys)
	for _, m := range many {
		f := []string{}
		for _, key := range keys {
			f = append(f, fmt.Sprint(m[key]))
		}
		w.Write(f)
	}
	w.Flush()
}
