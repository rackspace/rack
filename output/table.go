package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// ListTable writes a table composed of keys as the header with values from many
func ListTable(writer io.Writer, many []map[string]interface{}, keys []string) {
	w := tabwriter.NewWriter(writer, 0, 8, 1, '\t', 0)
	// Write the header
	fmt.Fprintln(w, strings.Join(keys, "\t"))
	for _, m := range many {
		f := []string{}
		for _, key := range keys {
			f = append(f, fmt.Sprint(m[key]))
		}
		fmt.Fprintln(w, strings.Join(f, "\t"))
	}
	w.Flush()
}

// MetadataTable writes a table to the writer composed of keys on the left and
// the associated metadata on the right column from m
func MetadataTable(writer io.Writer, m map[string]interface{}, keys []string) {
	w := tabwriter.NewWriter(writer, 0, 8, 0, '\t', 0)
	for _, key := range keys {
		val := fmt.Sprint(m[key])
		fmt.Fprintf(w, "%s\t%s\n", key, strings.Replace(val, "\n", "\n\t", -1))
	}
	w.Flush()
}
