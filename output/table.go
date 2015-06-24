package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

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

func MetadataTable(writer io.Writer, m map[string]interface{}, keys []string) {
	w := tabwriter.NewWriter(writer, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "PROPERTY\tVALUE")
	for _, key := range keys {
		val := fmt.Sprint(m[key])
		fmt.Fprintf(w, "%s\t%s\n", key, strings.Replace(val, "\n", "\n\t", -1))
	}
	w.Flush()
}
