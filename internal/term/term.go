// Package term provides minimal ANSI-colored terminal output with automatic
// TTY and NO_COLOR detection. No third-party dependencies are used.
package term

import (
	"os"
	"strings"
)

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	gray   = "\033[90m"
	bold   = "\033[1m"
)

// Enabled controls whether color codes are emitted.
var Enabled = autoColor()

func autoColor() bool {
	if v := os.Getenv("NO_COLOR"); v != "" {
		return false
	}
	if v := strings.ToLower(os.Getenv("IGNOREKIT_COLOR")); v == "never" || v == "false" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func paint(code, s string) string {
	if !Enabled {
		return s
	}
	return code + s + reset
}

// Colorizers.
func Red(s string) string    { return paint(red, s) }
func Green(s string) string  { return paint(green, s) }
func Yellow(s string) string { return paint(yellow, s) }
func Blue(s string) string   { return paint(blue, s) }
func Purple(s string) string { return paint(purple, s) }
func Cyan(s string) string   { return paint(cyan, s) }
func Gray(s string) string   { return paint(gray, s) }
func Bold(s string) string   { return paint(bold, s) }
