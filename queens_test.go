package main

import "testing"

func TestQueensPlacement(t *testing.T) {
	q := NewQueens()

	// Should be able to place first queen
	if err := q.PlaceQueen(0, 0); err != nil {
		t.Errorf("Failed to place first queen: %v", err)
	}

	// Should not be able to place on same row
	if err := q.PlaceQueen(0, 5); err == nil {
		t.Errorf("Should not allow placing on same row")
	}

	// Should not be able to place on same column
	if err := q.PlaceQueen(5, 0); err == nil {
		t.Errorf("Should not allow placing on same column")
	}

	// Should not be able to place on diagonal
	if err := q.PlaceQueen(1, 1); err == nil {
		t.Errorf("Should not allow placing on diagonal")
	}

	// Should be able to place on safe position
	if err := q.PlaceQueen(1, 2); err != nil {
		t.Errorf("Failed to place queen on safe position: %v", err)
	}

	if q.Count() != 2 {
		t.Errorf("Expected 2 queens, got %d", q.Count())
	}
}

func TestQueensRemoval(t *testing.T) {
	q := NewQueens()

	// Place a queen
	q.PlaceQueen(0, 0)

	// Remove it
	if err := q.RemoveQueen(0, 0); err != nil {
		t.Errorf("Failed to remove queen: %v", err)
	}

	if q.Count() != 0 {
		t.Errorf("Expected 0 queens after removal, got %d", q.Count())
	}

	// Try to remove from empty position
	if err := q.RemoveQueen(0, 0); err == nil {
		t.Errorf("Should error when removing from empty position")
	}
}

func TestQueensAttackDetection(t *testing.T) {
	q := NewQueens()
	q.PlaceQueen(3, 3)

	tests := []struct {
		row      int
		col      int
		expected bool
	}{
		{3, 0, true},  // same row
		{3, 7, true},  // same row
		{0, 3, true},  // same column
		{7, 3, true},  // same column
		{0, 0, true},  // diagonal
		{6, 6, true},  // diagonal
		{1, 5, true},  // diagonal
		{5, 1, true},  // diagonal
		{0, 1, false}, // safe
		{2, 5, false}, // safe
	}

	for _, tt := range tests {
		result := q.IsUnderAttack(tt.row, tt.col)
		if result != tt.expected {
			t.Errorf("IsUnderAttack(%d, %d) = %v, want %v", tt.row, tt.col, result, tt.expected)
		}
	}
}

func TestQueensSolution(t *testing.T) {
	q := NewQueens()

	// One valid solution for 8-queens
	solution := []Position{
		{0, 0},
		{1, 4},
		{2, 7},
		{3, 5},
		{4, 2},
		{5, 6},
		{6, 1},
		{7, 3},
	}

	for _, pos := range solution {
		if err := q.PlaceQueen(pos.Row, pos.Col); err != nil {
			t.Errorf("Failed to place queen at (%d, %d): %v", pos.Row, pos.Col, err)
		}
	}

	if !q.IsSolved() {
		t.Errorf("Expected puzzle to be solved")
	}

	if q.Count() != 8 {
		t.Errorf("Expected 8 queens, got %d", q.Count())
	}
}

func TestQueensSymbol(t *testing.T) {
	q := NewQueens()

	// Default should be black
	if q.GetSymbol() != "♛" {
		t.Errorf("Expected default symbol to be ♛, got %s", q.GetSymbol())
	}

	// Test changing symbols
	q.SetSymbol(SymbolWhite)
	if q.GetSymbol() != "♕" {
		t.Errorf("Expected white symbol ♕, got %s", q.GetSymbol())
	}

	q.SetSymbol(SymbolAscii)
	if q.GetSymbol() != "Q" {
		t.Errorf("Expected ASCII symbol Q, got %s", q.GetSymbol())
	}
}

func TestQueensReset(t *testing.T) {
	q := NewQueens()

	// Place some queens
	q.PlaceQueen(0, 0)
	q.PlaceQueen(1, 2)
	q.PlaceQueen(2, 4)

	if q.Count() != 3 {
		t.Errorf("Expected 3 queens before reset, got %d", q.Count())
	}

	// Reset
	q.Reset()

	if q.Count() != 0 {
		t.Errorf("Expected 0 queens after reset, got %d", q.Count())
	}

	// Should be able to place queens again
	if err := q.PlaceQueen(0, 0); err != nil {
		t.Errorf("Failed to place queen after reset: %v", err)
	}
}
