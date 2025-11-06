package main

import (
	"errors"
	"fmt"
	"strings"
)

const (
	boardSize = 8
)

var (
	ErrOutOfBounds     = errors.New("out of bounds")
	ErrOccupied        = errors.New("cell is occupied")
	ErrUnderAttack     = errors.New("cell is under attack")
	ErrNoQueenToRemove = errors.New("no queen to remove")
)

type Position struct {
	Row int
	Col int
}

type QueenSymbol int

const (
	SymbolBlack QueenSymbol = iota
	SymbolWhite
	SymbolAscii
)

type Queens struct {
	queens []Position
	symbol QueenSymbol
}

func NewQueens() Queens {
	return Queens{
		queens: make([]Position, 0, boardSize),
		symbol: SymbolBlack,
	}
}

func (q *Queens) SetSymbol(symbol QueenSymbol) {
	q.symbol = symbol
}

func (q *Queens) GetSymbol() string {
	switch q.symbol {
	case SymbolBlack:
		return "♛"
	case SymbolWhite:
		return "♕"
	case SymbolAscii:
		return "Q"
	default:
		return "Q"
	}
}

func (q *Queens) PlaceQueen(row, col int) error {
	if !inBounds(row, col) {
		return ErrOutOfBounds
	}

	if q.HasQueen(row, col) {
		return ErrOccupied
	}

	if q.IsUnderAttack(row, col) {
		return ErrUnderAttack
	}

	q.queens = append(q.queens, Position{Row: row, Col: col})
	return nil
}

func (q *Queens) RemoveQueen(row, col int) error {
	if !inBounds(row, col) {
		return ErrOutOfBounds
	}

	for i, pos := range q.queens {
		if pos.Row == row && pos.Col == col {
			// Remove queen at index i
			q.queens = append(q.queens[:i], q.queens[i+1:]...)
			return nil
		}
	}

	return ErrNoQueenToRemove
}

func (q *Queens) HasQueen(row, col int) bool {
	for _, pos := range q.queens {
		if pos.Row == row && pos.Col == col {
			return true
		}
	}
	return false
}

func (q *Queens) IsUnderAttack(row, col int) bool {
	for _, queen := range q.queens {
		if queen.Row == row {
			return true
		}
		if queen.Col == col {
			return true
		}
		rowDiff := abs(queen.Row - row)
		colDiff := abs(queen.Col - col)
		if rowDiff == colDiff {
			return true
		}
	}
	return false
}

// IsQueenUnderAttack checks if a queen at the given position is under attack by any OTHER queen
func (q *Queens) IsQueenUnderAttack(row, col int) bool {
	for _, queen := range q.queens {
		if queen.Row == row && queen.Col == col {
			continue
		}

		if queen.Row == row {
			return true
		}
		if queen.Col == col {
			return true
		}
		rowDiff := abs(queen.Row - row)
		colDiff := abs(queen.Col - col)
		if rowDiff == colDiff {
			return true
		}
	}
	return false
}

func (q *Queens) GetAttackedPositions() map[Position]bool {
	attacked := make(map[Position]bool)

	for _, queen := range q.queens {
		for col := 0; col < boardSize; col++ {
			if col != queen.Col {
				attacked[Position{Row: queen.Row, Col: col}] = true
			}
		}

		for row := 0; row < boardSize; row++ {
			if row != queen.Row {
				attacked[Position{Row: row, Col: queen.Col}] = true
			}
		}

		for i := 1; i < boardSize; i++ {
			if inBounds(queen.Row-i, queen.Col-i) {
				attacked[Position{Row: queen.Row - i, Col: queen.Col - i}] = true
			}
			if inBounds(queen.Row-i, queen.Col+i) {
				attacked[Position{Row: queen.Row - i, Col: queen.Col + i}] = true
			}
			if inBounds(queen.Row+i, queen.Col-i) {
				attacked[Position{Row: queen.Row + i, Col: queen.Col - i}] = true
			}
			if inBounds(queen.Row+i, queen.Col+i) {
				attacked[Position{Row: queen.Row + i, Col: queen.Col + i}] = true
			}
		}
	}

	return attacked
}

func (q *Queens) Count() int {
	return len(q.queens)
}

func (q *Queens) IsSolved() bool {
	if len(q.queens) != boardSize {
		return false
	}

	for _, queen := range q.queens {
		if q.IsQueenUnderAttack(queen.Row, queen.Col) {
			return false
		}
	}

	return true
}

func (q *Queens) Reset() {
	q.queens = make([]Position, 0, boardSize)
}

func (q *Queens) Pretty(cursorRow, cursorCol int, showAttacked bool, hardMode bool) string {
	var result strings.Builder

	attacked := make(map[Position]bool)
	if showAttacked && !hardMode {
		attacked = q.GetAttackedPositions()
	}

	queenSymbol := q.GetSymbol()

	result.WriteString("┌")
	for col := 0; col < boardSize; col++ {
		result.WriteString("───")
		if col < boardSize-1 {
			result.WriteString("┬")
		}
	}
	result.WriteString("┐\n")

	for row := 0; row < boardSize; row++ {
		result.WriteString("│")

		for col := 0; col < boardSize; col++ {
			isCursor := (row == cursorRow && col == cursorCol)
			hasQueen := q.HasQueen(row, col)
			isAttacked := attacked[Position{Row: row, Col: col}]

			if hasQueen {
				if hardMode {
					queenUnderAttack := q.IsQueenUnderAttack(row, col)
					if isCursor {
						if queenUnderAttack {
							result.WriteString(fmt.Sprintf("\033[1;31;7m %s \033[0m", queenSymbol))
						} else {
							result.WriteString(fmt.Sprintf("\033[1;32;7m %s \033[0m", queenSymbol))
						}
					} else {
						if queenUnderAttack {
							result.WriteString(fmt.Sprintf("\033[1;31m %s \033[0m", queenSymbol))
						} else {
							result.WriteString(fmt.Sprintf("\033[1;32m %s \033[0m", queenSymbol))
						}
					}
				} else {
					if isCursor {
						result.WriteString(fmt.Sprintf("\033[1;7m %s \033[0m", queenSymbol))
					} else {
						result.WriteString(fmt.Sprintf("\033[1m %s \033[0m", queenSymbol))
					}
				}
			} else if isCursor {
				result.WriteString("\033[1;7m   \033[0m")
			} else if isAttacked {
				result.WriteString("\033[41m   \033[0m")
			} else {
				result.WriteString("   ")
			}

			result.WriteString("│")
		}

		result.WriteString("\n")

		if row < boardSize-1 {
			result.WriteString("├")
			for col := 0; col < boardSize; col++ {
				result.WriteString("───")
				if col < boardSize-1 {
					result.WriteString("┼")
				}
			}
			result.WriteString("┤\n")
		}
	}

	result.WriteString("└")
	for col := 0; col < boardSize; col++ {
		result.WriteString("───")
		if col < boardSize-1 {
			result.WriteString("┴")
		}
	}
	result.WriteString("┘")

	return result.String()
}

func inBounds(row, col int) bool {
	return row >= 0 && row < boardSize && col >= 0 && col < boardSize
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
