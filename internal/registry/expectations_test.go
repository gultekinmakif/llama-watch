package registry

import "testing"

func TestClassifyCell_PresentMapsToPresent(t *testing.T) {
	if got := ClassifyCell(true); got != CellPresent {
		t.Fatalf("present cell: got %q, want %q", got, CellPresent)
	}
}

func TestClassifyCell_AbsentMapsToNA(t *testing.T) {
	if got := ClassifyCell(false); got != CellNA {
		t.Fatalf("absent cell: got %q, want %q", got, CellNA)
	}
}
