package output

import (
	"encoding/csv"
	"fmt"
	"io"
)

func MetadataCSV(writer io.Writer, m map[string]interface{}, keys []string, noHeader bool) {
	w := csv.NewWriter(writer)
	if !noHeader {
		w.Write([]string{"PROPERTY", "VALUE"})
	}
	for _, key := range keys {
		w.Write([]string{key, fmt.Sprint(m[key])})
	}
	w.Flush()
}

func ListCSV(writer io.Writer, many []map[string]interface{}, keys []string, noHeader bool) {
	w := csv.NewWriter(writer)
	if !noHeader {
		w.Write(keys)
	}
	for _, m := range many {
		f := []string{}
		for _, key := range keys {
			f = append(f, fmt.Sprint(m[key]))
		}
		w.Write(f)
	}
	w.Flush()
}
