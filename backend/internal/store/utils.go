package store

import (
	"fmt"
	"strings"
)

// buildBatchPlaceholders generates PostgreSQL placeholder groups for batch inserts.
// e.g. buildBatchPlaceholders(2, 3) → ["($1,$2,$3)", "($4,$5,$6)"]
func buildBatchPlaceholders(numRows, numCols int) []string {
	placeholders := make([]string, numRows)
	for i := range numRows {
		cols := make([]string, numCols)
		base := i * numCols
		for j := range numCols {
			cols[j] = fmt.Sprintf("$%d", base+j+1)
		}
		placeholders[i] = fmt.Sprintf("(%s)", strings.Join(cols, ","))
	}
	return placeholders
}
