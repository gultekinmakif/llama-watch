// Two-state cell classifier. Web client owns the four-state coloring.
package registry

type CellState string

const (
	CellNA      CellState = "na"
	CellPresent CellState = "present"
)

func ClassifyCell(present bool) CellState {
	if present {
		return CellPresent
	}
	return CellNA
}
