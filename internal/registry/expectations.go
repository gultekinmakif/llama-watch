// Two-state cell classifier. Web client owns the four-state coloring.
package registry

type CellState string

const (
	CellNA      CellState = "na"
	CellMissing CellState = "missing"
	CellPresent CellState = "present"
	CellOver    CellState = "over"
)

func ClassifyCell(category, metric string, present bool) CellState {
	_ = category
	_ = metric
	if present {
		return CellPresent
	}
	return CellNA
}
