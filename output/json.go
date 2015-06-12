package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSON prints results in JSON format.
func defaultJSON(w io.Writer, i interface{}) {
	m := map[string]interface{}{"result": i}
	j, _ := json.Marshal(m)
	fmt.Fprintf(w, "%v", string(j))
}

func metadataJSON(w io.Writer, m map[string]interface{}, keys []string) {
	mLimited := limitJSONFields(m, keys)
	j, _ := json.Marshal(mLimited)
	fmt.Fprintf(w, "%v", string(j))
}

func listJSON(w io.Writer, maps []map[string]interface{}, keys []string) {
	mLimited := make([]map[string]interface{}, len(maps))
	for i, m := range maps {
		mLimited[i] = limitJSONFields(m, keys)
	}
	j, _ := json.Marshal(mLimited)
	fmt.Fprintf(w, "%v", string(j))
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
