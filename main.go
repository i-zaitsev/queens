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

func checkAndUpdateSolution(queens Queens, userConfig *UserConfig, fundamentals [][]Position) {
	if queens.IsSolved() {
		matchNum := FindMatchingSolution(queens.queens, fundamentals)
		if matchNum != -1 && userConfig.Solved[matchNum-1] == 0 {
			userConfig.MarkSolved(matchNum)
			SaveConfig(userConfig)
		}
	}
}

func countSolved(solved [12]int) int {
	count := 0
	for _, s := range solved {
		count += s
	}
	return count
}

func main() {
	noExit := flag.Bool("noexit", false, "disable Esc; use :q to exit")
	hard := flag.Bool("hard", false, "hard mode: no help, show queen validity")
	flag.Parse()

	fundamentalSolutions, err := LoadFundamentalSolutions()
	if err != nil {
		panic(fmt.Errorf("failed to load fundamental solutions: %v", err))
	}

	userConfig, err := LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load user config: %v", err))
	}

	prizes, err := LoadPrizes()
	if err != nil {
		panic(fmt.Errorf("failed to load prizes: %v", err))
	}

	terminal := RawTerminal(*noExit)
	defer terminal.Restore()

	enterAltScreen()
	defer exitAltScreen()

	queens := NewQueens()
	cursorRow, cursorCol := 0, 0
	showHelp := false
	commandMode := false
	commandBuffer := ""

	renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)

	for {
		if cmd, err := terminal.ReadInput(); err == nil {
			if commandMode {
				switch cmd.Code {
				case CodeExit, CodeCancelCommand:
					commandMode = false
					terminal.SetCommandMode(false)
					commandBuffer = ""
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
				case CodePlace:
					if commandBuffer == ":q" {
						return
					}
					commandMode = false
					terminal.SetCommandMode(false)
					commandBuffer = ""
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
				case CodeChar:
					if data, ok := cmd.Data.(rune); ok {
						commandBuffer += string(data)
						renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
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
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)

			case CodeReset:
				queens.Reset()
				cursorRow, cursorCol = 0, 0
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)

			case CodeHelp:
				if !*hard {
					showHelp = !showHelp
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
				}

			case CodeSymbolBlack:
				queens.SetSymbol(SymbolBlack)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)

			case CodeSymbolWhite:
				queens.SetSymbol(SymbolWhite)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)

			case CodeSymbolAscii:
				queens.SetSymbol(SymbolAscii)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)

			case CodePlace:
				if queens.HasQueen(cursorRow, cursorCol) {
					queens.RemoveQueen(cursorRow, cursorCol)
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
				} else if queens.Count() < 8 {
					if *hard {
						queens.queens = append(queens.queens, Position{Row: cursorRow, Col: cursorCol})
						checkAndUpdateSolution(queens, userConfig, fundamentalSolutions)
						renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
					} else {
						if err := queens.PlaceQueen(cursorRow, cursorCol); err == nil {
							checkAndUpdateSolution(queens, userConfig, fundamentalSolutions)
							renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
						} else {
							renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
						}
					}
				}

			case CodeUp:
				if cursorRow > 0 {
					cursorRow--
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
				}

			case CodeDown:
				if cursorRow < boardSize-1 {
					cursorRow++
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
				}

			case CodeLeft:
				if cursorCol > 0 {
					cursorCol--
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
				}

			case CodeRight:
				if cursorCol < boardSize-1 {
					cursorCol++
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, userConfig.Solved, prizes)
				}

			case CodeNone:
			}
		} else {
			panic("error reading from terminal")
		}
	}
}
