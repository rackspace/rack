package output

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
)

// ListTable writes a table listing from an array of map[string]interface{}
func ListTable(c *cli.Context, many []map[string]interface{}, keys []string) {
	w := tabwriter.NewWriter(c.App.Writer, 0, 8, 1, '\t', 0)
	// Write the header
	fmt.Fprintln(w, strings.Join(keys, "\t"))

	for _, m := range many {
		writeMapEntry(w, m, keys)
	}
	w.Flush()
}

// MetaDataTable writes standardized metadata out
// Solely a utility method that invokes MetaDataMapTable with structs.Map(i)
func MetaDataTable(c *cli.Context, i interface{}, keys []string) {
	m := structs.Map(i)
	MetaDataMapTable(c, m, keys)
}

// MetaDataMapTable writes standardized metadata out
func MetaDataMapTable(c *cli.Context, m map[string]interface{}, keys []string) {
	w := tabwriter.NewWriter(c.App.Writer, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, "PROPERTY\tVALUE")

	for _, key := range keys {
		val := fmt.Sprint(m[key])
		fmt.Fprintf(w, "%s\t%s\n", key, strings.Replace(val, "\n", "\n\t", -1))
	}
	w.Flush()
}

// writeMapEntry writes a table entry from a map
func writeMapEntry(w *tabwriter.Writer, m map[string]interface{}, keys []string) {
	f := []string{}
	for _, key := range keys {
		f = append(f, fmt.Sprint(m[key]))
	}
	fmt.Fprintln(w, strings.Join(f, "\t"))
}
