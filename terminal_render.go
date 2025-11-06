package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func renderScreen(queens Queens, cursorRow, cursorCol int, showHelp bool, noExit bool, hard bool, commandBuffer string, solved [12]int) {
	fmt.Print("\033[H\033[2J")

	termWidth := getTerminalWidth()
	isSolved := queens.IsSolved()

	renderTitle(termWidth, isSolved)

	prettyString := queens.Pretty(cursorRow, cursorCol, showHelp, hard)
	lines := strings.Split(prettyString, "\n")

	for _, line := range lines {
		printCentered(line, termWidth)
	}

	fmt.Print("\r\n")

	renderStatus(queens, showHelp, termWidth, isSolved, hard)

	fmt.Print("\r\n")

	renderDiscoveryGrid(termWidth, solved)

	renderControls(termWidth, isSolved, noExit, hard)

	if commandBuffer != "" {
		renderCommandLine(commandBuffer, termWidth)
	}
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}
	return width
}

func printCentered(line string, termWidth int) {
	visibleLen := getVisibleLength(line)

	if visibleLen >= termWidth {
		fmt.Print(line)
		fmt.Print("\r\n")
		return
	}

	leftPadding := (termWidth - visibleLen) / 2

	fmt.Print(strings.Repeat(" ", leftPadding))
	fmt.Print(line)
	fmt.Print("\r\n")
}

func getVisibleLength(s string) int {
	length := 0
	inEscape := false

	for _, char := range s {
		if char == '\033' {
			inEscape = true
		} else if inEscape {
			if char == 'm' {
				inEscape = false
			}
		} else {
			length++
		}
	}

	return length
}

func renderTitle(termWidth int, isSolved bool) {
	fmt.Print("\033[33m")
	printCentered("╔════════════════════════════╗", termWidth)
	printCentered("║   8-Queens Puzzle (v1.0)   ║", termWidth)
	printCentered("╚════════════════════════════╝", termWidth)
	fmt.Print("\033[0m")
	fmt.Print("\r\n")
}

func renderDiscoveryGrid(termWidth int, solved [12]int) {
	fmt.Print("\033[36m")
	printCentered("Fundamental Solutions:", termWidth)
	fmt.Print("\033[0m")

	row1 := ""
	for i := 0; i < 6; i++ {
		cellNum := fmt.Sprintf("%02d", i+1)
		if solved[i] == 1 {
			row1 += fmt.Sprintf("\033[32m[%s]\033[0m ", cellNum)
		} else {
			row1 += fmt.Sprintf("[%s] ", cellNum)
		}
	}
	printCentered(strings.TrimSpace(row1), termWidth)

	row2 := ""
	for i := 6; i < 12; i++ {
		cellNum := fmt.Sprintf("%02d", i+1)
		if solved[i] == 1 {
			row2 += fmt.Sprintf("\033[32m[%s]\033[0m ", cellNum)
		} else {
			row2 += fmt.Sprintf("[%s] ", cellNum)
		}
	}
	printCentered(strings.TrimSpace(row2), termWidth)

	fmt.Print("\r\n")
}

func renderStatus(queens Queens, showHelp bool, termWidth int, isSolved bool, hard bool) {
	fmt.Print("\033[32m")

	status := fmt.Sprintf("Queens: %d/8", queens.Count())
	if isSolved {
		status += "  \033[1;32m✓ Solved!\033[0m\033[32m"
	}

	status += fmt.Sprintf("  Symbol: %s", queens.GetSymbol())

	if !hard {
		helpStatus := "OFF"
		if showHelp {
			helpStatus = "ON"
		}
		status += fmt.Sprintf("  Help: %s", helpStatus)
	}

	printCentered(status, termWidth)
	fmt.Print("\033[0m")
	fmt.Print("\r\n")
}

func renderControls(termWidth int, isSolved bool, noExit bool, hard bool) {
	fmt.Print("\033[36m")
	printCentered("┌────────────────────────────┐", termWidth)
	printCentered("│ Controls:                  │", termWidth)
	if !noExit {
		printCentered("│ [Esc]       Exit           │", termWidth)
	}
	printCentered("│ [r]         Reset board    │", termWidth)
	if !hard {
		printCentered("│ [h]         Toggle help    │", termWidth)
	}
	printCentered("│ [Space]     Toggle queen   │", termWidth)
	printCentered("│ [b/w/q]     Change symbol  │", termWidth)
	printCentered("│ [Arrows]    Move cursor    │", termWidth)
	printCentered("└────────────────────────────┘", termWidth)
	fmt.Print("\033[0m")
}

func renderCommandLine(commandBuffer string, termWidth int) {
	fmt.Print("\r\n")
	fmt.Print("\033[33m")
	printCentered(commandBuffer, termWidth)
	fmt.Print("\033[0m")
}
