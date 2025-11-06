package main

import (
	"flag"
	"fmt"
)

func enterAltScreen() {
	fmt.Print("\033[?1049h") // Enter alt screen buffer
	fmt.Print("\033[?25l")   // Hide cursor
}

func exitAltScreen() {
	fmt.Print("\033[?25h")   // Show cursor
	fmt.Print("\033[?1049l") // Exit alt screen buffer
}

func main() {
	noExit := flag.Bool("noexit", false, "disable Esc; use :q to exit")
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

	renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)

	for {
		if cmd, err := terminal.ReadInput(); err == nil {
			// Handle command mode separately
			if commandMode {
				switch cmd.Code {
				case CodeExit, CodeCancelCommand:
					// Esc or cancel in command mode
					commandMode = false
					terminal.SetCommandMode(false)
					commandBuffer = ""
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
				case CodePlace:
					// Enter pressed - evaluate command
					if commandBuffer == ":q" {
						return
					}
					// Invalid command, clear and exit command mode
					commandMode = false
					terminal.SetCommandMode(false)
					commandBuffer = ""
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
				case CodeChar:
					// Add character to command buffer
					if data, ok := cmd.Data.(rune); ok {
						commandBuffer += string(data)
						renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
					}
				}
				continue
			}

			// Normal mode
			switch cmd.Code {
			case CodeExit:
				return

			case CodeCommand:
				// Start command mode
				commandMode = true
				terminal.SetCommandMode(true)
				commandBuffer = ":"
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)

			case CodeReset:
				queens.Reset()
				cursorRow, cursorCol = 0, 0
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)

			case CodeHelp:
				showHelp = !showHelp
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)

			case CodeSymbolBlack:
				queens.SetSymbol(SymbolBlack)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)

			case CodeSymbolWhite:
				queens.SetSymbol(SymbolWhite)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)

			case CodeSymbolAscii:
				queens.SetSymbol(SymbolAscii)
				renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)

			case CodePlace:
				// Try to place a queen at the current cursor position
				if err := queens.PlaceQueen(cursorRow, cursorCol); err == nil {
					// Successfully placed, render
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
				} else {
					// Failed to place (occupied or under attack), just re-render
					// The visual feedback (red cells) should indicate why
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
				}

			case CodeRemove:
				// Try to remove a queen at the current cursor position
				if err := queens.RemoveQueen(cursorRow, cursorCol); err == nil {
					// Successfully removed
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
				}
				// If no queen to remove, just ignore silently

			case CodeUp:
				if cursorRow > 0 {
					cursorRow--
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
				}

			case CodeDown:
				if cursorRow < boardSize-1 {
					cursorRow++
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
				}

			case CodeLeft:
				if cursorCol > 0 {
					cursorCol--
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
				}

			case CodeRight:
				if cursorCol < boardSize-1 {
					cursorCol++
					renderScreen(queens, cursorRow, cursorCol, showHelp, *noExit, commandBuffer)
				}

			case CodeNone:
				// Ignore unknown input
			}
		} else {
			panic("error reading from terminal")
		}
	}
}
