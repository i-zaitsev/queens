package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Transform int

const (
	TransformIdentity Transform = iota
	TransformRot90
	TransformRot180
	TransformRot270
	TransformMirrorH
	TransformMirrorV
	TransformMirrorD
	TransformMirrorAD
)

type Symmetries struct {
	transforms map[Transform][]Position
}

func NewSymmetries(positions []Position) Symmetries {
	s := Symmetries{
		transforms: make(map[Transform][]Position),
	}

	s.transforms[TransformIdentity] = normalizePositions(positions)
	s.transforms[TransformRot90] = normalizePositions(rotate90(positions))
	s.transforms[TransformRot180] = normalizePositions(rotate180(positions))
	s.transforms[TransformRot270] = normalizePositions(rotate270(positions))
	s.transforms[TransformMirrorH] = normalizePositions(mirrorHorizontal(positions))
	s.transforms[TransformMirrorV] = normalizePositions(mirrorVertical(positions))
	s.transforms[TransformMirrorD] = normalizePositions(mirrorDiagonal(positions))
	s.transforms[TransformMirrorAD] = normalizePositions(mirrorAntiDiagonal(positions))

	return s
}

func (s *Symmetries) Matches(board []Position) bool {
	normalized := normalizePositions(board)
	for _, transform := range s.transforms {
		if positionsEqual(normalized, transform) {
			return true
		}
	}
	return false
}

func rotate90(positions []Position) []Position {
	result := make([]Position, len(positions))
	for i, pos := range positions {
		result[i] = Position{Row: pos.Col, Col: boardSize - 1 - pos.Row}
	}
	return result
}

func rotate180(positions []Position) []Position {
	result := make([]Position, len(positions))
	for i, pos := range positions {
		result[i] = Position{Row: boardSize - 1 - pos.Row, Col: boardSize - 1 - pos.Col}
	}
	return result
}

func rotate270(positions []Position) []Position {
	result := make([]Position, len(positions))
	for i, pos := range positions {
		result[i] = Position{Row: boardSize - 1 - pos.Col, Col: pos.Row}
	}
	return result
}

func mirrorHorizontal(positions []Position) []Position {
	result := make([]Position, len(positions))
	for i, pos := range positions {
		result[i] = Position{Row: pos.Row, Col: boardSize - 1 - pos.Col}
	}
	return result
}

func mirrorVertical(positions []Position) []Position {
	result := make([]Position, len(positions))
	for i, pos := range positions {
		result[i] = Position{Row: boardSize - 1 - pos.Row, Col: pos.Col}
	}
	return result
}

func mirrorDiagonal(positions []Position) []Position {
	result := make([]Position, len(positions))
	for i, pos := range positions {
		result[i] = Position{Row: pos.Col, Col: pos.Row}
	}
	return result
}

func mirrorAntiDiagonal(positions []Position) []Position {
	result := make([]Position, len(positions))
	for i, pos := range positions {
		result[i] = Position{Row: boardSize - 1 - pos.Col, Col: boardSize - 1 - pos.Row}
	}
	return result
}

func normalizePositions(positions []Position) []Position {
	result := make([]Position, len(positions))
	copy(result, positions)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Row != result[j].Row {
			return result[i].Row < result[j].Row
		}
		return result[i].Col < result[j].Col
	})
	return result
}

func positionsEqual(a, b []Position) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Row != b[i].Row || a[i].Col != b[i].Col {
			return false
		}
	}
	return true
}

func LoadFundamentalSolutions() ([][]Position, error) {
	files, err := filepath.Glob("boards/*.txt")
	if err != nil {
		return nil, err
	}

	sort.Strings(files)

	var solutions [][]Position
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		var positions []Position
		for row, line := range lines {
			line = strings.TrimSpace(line)
			col := strings.Index(line, "Q")
			if col != -1 {
				positions = append(positions, Position{Row: row, Col: col})
			}
		}
		solutions = append(solutions, positions)
	}

	return solutions, nil
}

func FindMatchingSolution(userBoard []Position, fundamentals [][]Position) int {
	for i, fundamental := range fundamentals {
		symmetries := NewSymmetries(fundamental)
		if symmetries.Matches(userBoard) {
			return i + 1
		}
	}
	return -1
}
