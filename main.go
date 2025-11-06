package main

import (
	"flag"
	"fmt"
)

func enterAltScreen() {
	fmt.Print("\033[?1049h")
	fmt.Print("\033[?25l")
}

func exitAltScreen() {
	fmt.Print("\033[?25h")
	fmt.Print("\033[?1049l")
}

func main() {
	noExit := flag.Bool("noexit", false, "disable Esc; use :q to exit")
	hard := flag.Bool("hard", false, "hard mode: no help, show queen validity")
	flag.Parse()

	terminal := RawTerminal(*noExit)
	defer terminal.Restore()

	enterAltScreen()
	defer exitAltScreen()

	queens := NewQueens()
	cursorRow, cursorCol := 0, 0
	showHelp := false
	commandMode := false
	commandBuffer := ""

	renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)

	for {
		if cmd, err := terminal.ReadInput(); err == nil {
			if commandMode {
				switch cmd.Code {
				case CodeExit, CodeCancelCommand:
					commandMode = false
					terminal.SetCommandMode(false)
					commandBuffer = ""
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
				case CodePlace:
					if commandBuffer == ":q" {
						return
					}
					commandMode = false
					terminal.SetCommandMode(false)
					commandBuffer = ""
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
				case CodeChar:
					if data, ok := cmd.Data.(rune); ok {
						commandBuffer += string(data)
						renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
					}
				}
				continue
			}

			switch cmd.Code {
			case CodeExit:
				return

			case CodeCommand:
				commandMode = true
				terminal.SetCommandMode(true)
				commandBuffer = ":"
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)

			case CodeReset:
				queens.Reset()
				cursorRow, cursorCol = 0, 0
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)

			case CodeHelp:
				if !*hard {
					showHelp = !showHelp
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
				}

			case CodeSymbolBlack:
				queens.SetSymbol(SymbolBlack)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)

			case CodeSymbolWhite:
				queens.SetSymbol(SymbolWhite)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)

			case CodeSymbolAscii:
				queens.SetSymbol(SymbolAscii)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)

			case CodePlace:
				if queens.HasQueen(cursorRow, cursorCol) {
					queens.RemoveQueen(cursorRow, cursorCol)
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
				} else if queens.Count() < 8 {
					if *hard {
						queens.queens = append(queens.queens, Position{Row: cursorRow, Col: cursorCol})
						renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
					} else {
						if err := queens.PlaceQueen(cursorRow, cursorCol); err == nil {
							renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
						} else {
							renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
						}
					}
				}

			case CodeUp:
				if cursorRow > 0 {
					cursorRow--
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
				}

			case CodeDown:
				if cursorRow < boardSize-1 {
					cursorRow++
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
				}

			case CodeLeft:
				if cursorCol > 0 {
					cursorCol--
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
				}

			case CodeRight:
				if cursorCol < boardSize-1 {
					cursorCol++
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer)
				}

			case CodeNone:
			}
		} else {
			panic("error reading from terminal")
		}
	}
}
