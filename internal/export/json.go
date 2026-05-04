package export

import (
	"encoding/json"

	"devmemory/internal/core"
)

// ExportJSON serializes entries to pretty-printed JSON.
func ExportJSON(entries []*core.Entry) ([]byte, error) {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ImportJSON deserializes entries from JSON.
func ImportJSON(data []byte) ([]*core.Entry, error) {
	var entries []*core.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
