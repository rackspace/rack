package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// INDENT is the indentation passed to json.MarshalIndent
const INDENT string = "  "

// JSON prints results in JSON format.
func defaultJSON(w io.Writer, i interface{}) {
	m := map[string]interface{}{"result": i}
	jsonOut(w, m)
}

func metadataJSON(w io.Writer, m map[string]interface{}, keys []string) {
	mLimited := limitJSONFields(m, keys)
	jsonOut(w, mLimited)
}

func listJSON(w io.Writer, maps []map[string]interface{}, keys []string) {
	mLimited := make([]map[string]interface{}, len(maps))
	for i, m := range maps {
		mLimited[i] = limitJSONFields(m, keys)
	}
	jsonOut(w, mLimited)
}

func limitJSONFields(m map[string]interface{}, keys []string) map[string]interface{} {
	mLimited := make(map[string]interface{})
	for _, key := range keys {
		if v, ok := m[key]; ok {
			mLimited[key] = v
		}
	}
	return mLimited
}

func jsonOut(w io.Writer, i interface{}) {
	j, _ := json.MarshalIndent(i, "", INDENT)
	fmt.Fprintln(w, string(j))
}
