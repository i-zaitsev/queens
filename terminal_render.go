package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func renderScreen(queens Queens, cursorRow, cursorCol int, showHelp bool, noExit bool, hard bool, commandBuffer string) {
	// Clear screen and move cursor to top-left
	fmt.Print("\033[H\033[2J")

	// Get terminal width for centering
	termWidth := getTerminalWidth()
	isSolved := queens.IsSolved()

	renderTitle(termWidth, isSolved)

	// Render the board
	prettyString := queens.Pretty(cursorRow, cursorCol, showHelp, hard)
	lines := strings.Split(prettyString, "\n")

	for _, line := range lines {
		printCentered(line, termWidth)
	}

	fmt.Print("\r\n")

	// Render status
	renderStatus(queens, showHelp, termWidth, isSolved, hard)

	// Render controls
	renderControls(termWidth, isSolved, noExit, hard)

	// Render command line if in command mode
	if commandBuffer != "" {
		renderCommandLine(commandBuffer, termWidth)
	}
}

// getTerminalWidth returns the terminal width, or 80 if unable to determine
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // Default fallback width
	}
	return width
}

// printCentered prints a line with left padding to center it horizontally
func printCentered(line string, termWidth int) {
	// Remove ANSI escape codes for accurate length calculation
	visibleLen := getVisibleLength(line)

	if visibleLen >= termWidth {
		// Line is too long, print without padding
		fmt.Print(line)
		fmt.Print("\r\n")
		return
	}

	// Calculate left padding
	leftPadding := (termWidth - visibleLen) / 2

	// Print padding and line
	fmt.Print(strings.Repeat(" ", leftPadding))
	fmt.Print(line)
	fmt.Print("\r\n")
}

// getVisibleLength calculates the visible length of a string, excluding ANSI escape codes
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
	fmt.Print("\033[33m") // Yellow color
	printCentered("╔════════════════════════════╗", termWidth)
	printCentered("║   8-Queens Puzzle (v1.0)   ║", termWidth)
	printCentered("╚════════════════════════════╝", termWidth)
	fmt.Print("\033[0m") // Reset color
	fmt.Print("\r\n")
}

func renderStatus(queens Queens, showHelp bool, termWidth int, isSolved bool, hard bool) {
	fmt.Print("\033[32m") // Green color

	// Show queen count
	status := fmt.Sprintf("Queens: %d/8", queens.Count())
	if isSolved {
		status += "  \033[1;32m✓ Solved!\033[0m\033[32m"
	}

	// Show current symbol
	status += fmt.Sprintf("  Symbol: %s", queens.GetSymbol())

	// Show help status (only if not in hard mode)
	if !hard {
		helpStatus := "OFF"
		if showHelp {
			helpStatus = "ON"
		}
		status += fmt.Sprintf("  Help: %s", helpStatus)
	}

	printCentered(status, termWidth)
	fmt.Print("\033[0m") // Reset color
	fmt.Print("\r\n")
}

func renderControls(termWidth int, isSolved bool, noExit bool, hard bool) {
	fmt.Print("\033[36m") // Cyan color
	printCentered("┌────────────────────────────┐", termWidth)
	printCentered("│ Controls:                  │", termWidth)
	if !noExit {
		printCentered("│ [Esc]       Exit           │", termWidth)
	}
	printCentered("│ [r]         Reset board    │", termWidth)
	if !hard {
		printCentered("│ [h]         Toggle help    │", termWidth)
	}
	printCentered("│ [Space]     Place queen    │", termWidth)
	printCentered("│ [x/Bksp]    Remove queen   │", termWidth)
	printCentered("│ [b/w/q]     Change symbol  │", termWidth)
	printCentered("│ [Arrows]    Move cursor    │", termWidth)
	printCentered("└────────────────────────────┘", termWidth)
	fmt.Print("\033[0m") // Reset color
}

func renderCommandLine(commandBuffer string, termWidth int) {
	fmt.Print("\r\n")
	fmt.Print("\033[33m") // Yellow color for command line
	printCentered(commandBuffer, termWidth)
	fmt.Print("\033[0m") // Reset color
}
