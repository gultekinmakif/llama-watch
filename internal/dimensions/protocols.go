// LoadProtocols reads the JSON file produced by tools/extract-protocols.ts
// and returns one RawProtocol across /data1.ts through /data6.ts.
package dimensions

import (
	"encoding/json"
	"os"
	"sort"
)

// RawProtocol is one entry from the extracted protocols JSON.
type RawProtocol struct {
	Name       string            `json:"name"`
	Category   string            `json:"category"`
	Chains     []string          `json:"chains"`
	Module     string            `json:"module"`
	Dimensions map[string]string `json:"dimensions"`
	DataFile   string            `json:"-"` // set during load from the JSON key
}

// LoadProtocols decodes the JSON and tags each entry's DataFile with its source key.
func LoadProtocols(jsonPath string) ([]RawProtocol, error) {
	f, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var byFile map[string][]RawProtocol
	if err := json.NewDecoder(f).Decode(&byFile); err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(byFile))
	for k := range byFile {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var out []RawProtocol
	for _, k := range keys {
		entries := byFile[k]
		for i := range entries {
			entries[i].DataFile = k + ".ts"
		}
		out = append(out, entries...)
	}
	return out, nil
}
