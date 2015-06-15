package output

import (
	"encoding/csv"
	"fmt"
	"io"
)

func metadataCSV(writer io.Writer, m map[string]interface{}, keys []string) {
	w := csv.NewWriter(writer)
	w.Write([]string{"PROPERTY", "VALUE"})
	for _, key := range keys {
		w.Write([]string{key, fmt.Sprint(m[key])})
	}
	w.Flush()
}

func listCSV(writer io.Writer, many []map[string]interface{}, keys []string) {
	w := csv.NewWriter(writer)
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
