// SQL query layer for /api/chains.
package api

import (
	"context"
	"unicode"

	"gorm.io/gorm"
)

const listChainsSQL = `
SELECT chain AS key, COUNT(*) AS protocol_count
FROM (SELECT unnest(chains) AS chain FROM protocol_identities) sub
GROUP BY chain
ORDER BY chain
`

// listChains returns one entry per distinct chain across protocols, ascending by key.
// Label is set in Go via titlecase. The chains.json override is intentionally deferred.
func listChains(ctx context.Context, db *gorm.DB) ([]ChainEntry, error) {
	rows := make([]ChainEntry, 0)
	if err := db.WithContext(ctx).Raw(listChainsSQL).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].Label = titlecase(rows[i].Key)
	}
	return rows, nil
}

// titlecase uppercases the first rune and leaves the rest unchanged.
func titlecase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
