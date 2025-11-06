package main

import (
	"flag"
	"fmt"
	"os"
)

func enterAltScreen() {
	fmt.Print("\033[?1049h")
	fmt.Print("\033[?25l")
}

func exitAltScreen() {
	fmt.Print("\033[?25h")
	fmt.Print("\033[?1049l")
}

func checkAndUpdateSolution(queens Queens, config *Config, playerName string, fundamentals [][]Position) {
	if queens.IsSolved() {
		matchNum := FindMatchingSolution(queens.queens, fundamentals)
		playerData := GetPlayerData(config, playerName)
		if matchNum != -1 && playerData[matchNum-1] == 0 {
			playerData[matchNum-1] = 1
			SetPlayerData(config, playerName, playerData)
			SaveConfig(config)
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
	player := flag.String("player", "", "player name for tracking progress (required)")
	flag.Parse()

	if *player == "" {
		fmt.Println("Error: -player flag is required")
		fmt.Println("Usage: queens -player <name> [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fundamentalSolutions, err := LoadFundamentalSolutions()
	if err != nil {
		panic(fmt.Errorf("failed to load fundamental solutions: %v", err))
	}

	config, err := LoadConfig(*player)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %v", err))
	}

	prizes, err := LoadPrizes()
	if err != nil {
		panic(fmt.Errorf("failed to load prizes: %v", err))
	}

	playerSolved := GetPlayerData(config, *player)

	terminal := RawTerminal(*noExit)
	defer terminal.Restore()

	enterAltScreen()
	defer exitAltScreen()

	queens := NewQueens()
	cursorRow, cursorCol := 0, 0
	showHelp := false
	commandMode := false
	commandBuffer := ""

	renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)

	for {
		if cmd, err := terminal.ReadInput(); err == nil {
			if commandMode {
				switch cmd.Code {
				case CodeExit, CodeCancelCommand:
					commandMode = false
					terminal.SetCommandMode(false)
					commandBuffer = ""
					playerSolved = GetPlayerData(config, *player)
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
				case CodePlace:
					if commandBuffer == ":q" {
						return
					}
					commandMode = false
					terminal.SetCommandMode(false)
					commandBuffer = ""
					playerSolved = GetPlayerData(config, *player)
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
				case CodeChar:
					if data, ok := cmd.Data.(rune); ok {
						commandBuffer += string(data)
						renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
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
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)

			case CodeReset:
				queens.Reset()
				cursorRow, cursorCol = 0, 0
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)

			case CodeHelp:
				if !*hard {
					showHelp = !showHelp
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
				}

			case CodeSymbolBlack:
				queens.SetSymbol(SymbolBlack)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)

			case CodeSymbolWhite:
				queens.SetSymbol(SymbolWhite)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)

			case CodeSymbolAscii:
				queens.SetSymbol(SymbolAscii)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)

			case CodePlace:
				if queens.HasQueen(cursorRow, cursorCol) {
					queens.RemoveQueen(cursorRow, cursorCol)
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
				} else if queens.Count() < 8 {
					if *hard {
						queens.queens = append(queens.queens, Position{Row: cursorRow, Col: cursorCol})
						checkAndUpdateSolution(queens, config, *player, fundamentalSolutions)
						playerSolved = GetPlayerData(config, *player)
						renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
					} else {
						if err := queens.PlaceQueen(cursorRow, cursorCol); err == nil {
							checkAndUpdateSolution(queens, config, *player, fundamentalSolutions)
							playerSolved = GetPlayerData(config, *player)
							renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
						} else {
							renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
						}
					}
				}

			case CodeUp:
				if cursorRow > 0 {
					cursorRow--
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
				}

			case CodeDown:
				if cursorRow < boardSize-1 {
					cursorRow++
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
				}

			case CodeLeft:
				if cursorCol > 0 {
					cursorCol--
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
				}

			case CodeRight:
				if cursorCol < boardSize-1 {
					cursorCol++
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, *hard, commandBuffer, playerSolved, prizes)
				}

			case CodeNone:
			}
		} else {
			panic("error reading from terminal")
		}
	}
}
