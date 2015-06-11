package output

import (
	"encoding/csv"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

// ListTable writes a table listing from an array of map[string]interface{}
func ListTable(c *cli.Context, f *func() []map[string]interface{}, keys []string) {
	w := tabwriter.NewWriter(c.App.Writer, 0, 8, 1, '\t', 0)
	// Write the header
	fmt.Fprintln(w, strings.Join(keys, "\t"))

	many := (*f)()
	for _, m := range many {
		f := []string{}
		for _, key := range keys {
			f = append(f, fmt.Sprint(m[key]))
		}
		fmt.Fprintln(w, strings.Join(f, "\t"))
	}
	w.Flush()
}

// MetadataTable writes standardized metadata out
func MetadataTable(c *cli.Context, f *func() map[string]interface{}, keys []string) {
	if c.IsSet("csv") {
		w := csv.NewWriter(c.App.Writer)
		m := (*f)()
		w.Write([]string{"PROPERTY", "VALUE"})
		for _, key := range keys {
			val := fmt.Sprint(m[key])
			w.Write([]string{key, val})
		}
		w.Flush()
	} else {
		w := tabwriter.NewWriter(c.App.Writer, 0, 8, 0, '\t', 0)
		m := (*f)()
		fmt.Fprintln(w, "PROPERTY\tVALUE")
		for _, key := range keys {
			val := fmt.Sprint(m[key])
			fmt.Fprintf(w, "%s\t%s\n", key, strings.Replace(val, "\n", "\n\t", -1))
		}
		w.Flush()
	}
}
