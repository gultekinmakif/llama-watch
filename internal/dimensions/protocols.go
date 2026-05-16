// 1.3: Decodes the JSON file produced by tools/extract-protocols.ts.
// Second pipeline phase. Returns the map keyed by source data file (data1..data6).
package dimensions

import (
	"context"
	"encoding/json"
	"os"
)

// RawProtocol is one entry from the extracted protocols JSON.
type RawProtocol struct {
	Name       string            `json:"name"`
	Category   string            `json:"category"`
	Chains     []string          `json:"chains"`
	Module     string            `json:"module"`
	Dimensions map[string]string `json:"dimensions"`
}

// LoadProtocols decodes the extractor output. The map key is the source data file (e.g. "data1").
func LoadProtocols(ctx context.Context, jsonPath string) (map[string][]RawProtocol, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	f, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var byFile map[string][]RawProtocol
	err = json.NewDecoder(f).Decode(&byFile)
	return byFile, err
}
