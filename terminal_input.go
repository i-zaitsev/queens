package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	reset       func()
	noExit      bool
	commandMode bool
}

type Code int

const (
	CodeExit Code = iota
	CodeReset
	CodeLeft
	CodeRight
	CodeUp
	CodeDown
	CodePlace
	CodeRemove
	CodeHelp
	CodeSymbolBlack
	CodeSymbolWhite
	CodeSymbolAscii
	CodeCommand
	CodeCancelCommand
	CodeChar
	CodeNone
)

type Cmd struct {
	Code Code
	Data interface{}
}

func NewCmd(code Code) Cmd {
	return Cmd{
		Code: code,
	}
}

func RawTerminal(noExit bool) Terminal {
	state, err := term.MakeRaw(getTermFd())
	if err != nil {
		panic(fmt.Errorf("cannot set terminal to raw mode: %v", err))
	}
	return Terminal{
		reset: func() {
			if err := term.Restore(getTermFd(), state); err != nil {
				panic(fmt.Errorf("cannot restore terminal state: %v", err))
			}
		},
		noExit: noExit,
	}
}

func (t *Terminal) Restore() {
	t.reset()
}

func (t *Terminal) SetCommandMode(mode bool) {
	t.commandMode = mode
}

func (t *Terminal) ReadInput() (Cmd, error) {
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err == io.EOF {
		return NewCmd(CodeNone), nil
	} else if err != nil {
		return NewCmd(CodeNone), fmt.Errorf("error reading from stdin: %v", err)
	}
	if n > 0 {
		if n == 1 {
			char := buf[0]

			if t.commandMode {
				if char == 0x1b {
					return NewCmd(CodeCancelCommand), nil
				} else if char == '\r' || char == '\n' {
					return NewCmd(CodePlace), nil
				} else if isPrintable(char) {
					cmd := NewCmd(CodeChar)
					cmd.Data = rune(char)
					return cmd, nil
				}
				return NewCmd(CodeNone), nil
			}

			if char == 0x1b {
				if t.noExit {
					return NewCmd(CodeNone), nil
				}
				return NewCmd(CodeExit), nil
			} else if char == ':' {
				return NewCmd(CodeCommand), nil
			} else if char == 'r' || char == 'R' {
				return NewCmd(CodeReset), nil
			} else if char == 'h' || char == 'H' {
				return NewCmd(CodeHelp), nil
			} else if char == 'b' || char == 'B' {
				return NewCmd(CodeSymbolBlack), nil
			} else if char == 'w' || char == 'W' {
				return NewCmd(CodeSymbolWhite), nil
			} else if char == 'q' || char == 'Q' {
				return NewCmd(CodeSymbolAscii), nil
			} else if char == ' ' || char == '\r' || char == '\n' {
				return NewCmd(CodePlace), nil
			} else if char == 'x' || char == 'X' || char == 127 || char == 8 {
				return NewCmd(CodeRemove), nil
			}
		} else if n == 3 && buf[0] == 0x1b && buf[1] == '[' {
			switch buf[2] {
			case 'A':
				return NewCmd(CodeUp), nil
			case 'B':
				return NewCmd(CodeDown), nil
			case 'C':
				return NewCmd(CodeRight), nil
			case 'D':
				return NewCmd(CodeLeft), nil
			default:
				return NewCmd(CodeNone), nil
			}
		}
	}
	return NewCmd(CodeNone), nil
}

func isPrintable(char byte) bool {
	return char >= 32 && char <= 126
}

func getTermFd() int {
	return int(os.Stdin.Fd())
}
