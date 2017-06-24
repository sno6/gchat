package ui

import (
	"fmt"
	"io"
	"strconv"
)

// The Color interface defines anything that can take a string and return
// a colored string.
type Color interface {
	Color(string) string
}

type (
	// Normal represents a normal ANSI color code.
	Normal int

	// Bold is a bold ANSI color.
	Bold int
)

// Color a string normaly.
func (n Normal) Color(s string) string {
	return esc + "[" + strconv.Itoa(int(n)) + "m" + s + rst
}

// Color a string bold.
func (b Bold) Color(s string) string {
	return esc + "[" + strconv.Itoa(int(b)) + ";" + bold + s + rst
}

// ANSI color codes.
const (
	Red Normal = iota + 31
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	BoldRed     = Bold(Red)
	BoldGreen   = Bold(Green)
	BoldYellow  = Bold(Yellow)
	BoldBlue    = Bold(Blue)
	BoldMagenta = Bold(Magenta)
	BoldCyan    = Bold(Cyan)
	BoldWhite   = Bold(White)
)

const (
	esc  = "\033"
	rst  = esc + "[0m"
	bold = "1m"
)

// Sprintf decorates the fmt.Sprintf function to place ANSI color codes on a string.
func Sprintf(color Color, format string, a ...interface{}) string {
	return fmt.Sprintf(color.Color(format), a...)
}

// Fprintf decorates the fmt.Fprintf function to place ANSI color codes on a string
// which is writter to an io.Writer.
func Fprintf(w io.Writer, color Color, format string, a ...interface{}) {
	fmt.Fprintf(w, color.Color(format), a...)
}
