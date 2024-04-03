package color

import (
	"fmt"
	"runtime"
	"strings"
)

var Reset   = []byte("\033[0m")
var Black   = []byte("\033[30m")
var Red     = []byte("\033[31m")
var Green   = []byte("\033[32m")
var Yellow  = []byte("\033[33m")
var Blue    = []byte("\033[34m")
var Magenta = []byte("\033[35m")
var Cyan    = []byte("\033[36m")
var Gray    = []byte("\033[37m")
var White   = []byte("\033[97m")
var BrightBlack = []byte("\033[90m")
var BrightRed = []byte("\033[91m")
var BrightGreen = []byte("\033[92m")
var BrightYellow = []byte("\033[93m")
var BrightBlue = []byte("\033[94m")
var BrightMagenta = []byte("\033[95m")
var BrightCyan = []byte("\033[96m")

var Colors = [][]byte{Reset, 
Red, 
Green, 
Yellow, 
Blue, 
Magenta, 
Cyan, 
Gray, 
White,
BrightBlack,
BrightRed,
BrightGreen,
BrightYellow,
BrightBlue,
BrightMagenta,
BrightCyan,
}

func init() {
	if runtime.GOOS == "windows" {
        Colors = [][]byte{}
    }
}

// usage after importing package: colors.ColorPrintf(colors.Magenta, "My text here with a variable: %s", myVariable)
func ColorPrintf(color []byte, format string, a ...any) {
    formatBytes := []byte(format)

    if checkColor(color) {
        formatBytes = append(color, formatBytes...)
        formatBytes = append(formatBytes, Reset...)
    }

    // a is set to nil if not provided, and outputs garbage to the terminal.
    if a == nil {
        _, err := fmt.Printf(string(formatBytes))
        if err != nil {
            panic(err.Error())
        } 
    } else {
        _, err := fmt.Printf(string(formatBytes), a)
        if err != nil {
            panic(err.Error())
        }
    }
}

func ColorSprintf(color []byte, format string, a ...any) (s string) {
    formatBytes := []byte(format)

    if checkColor(color) {
        formatBytes = append(color, formatBytes...)
        formatBytes = append(formatBytes, Reset...)
    }

    // a is set to nil if not provided, and outputs garbage to the terminal.
    if a == nil {
        s = fmt.Sprintf(string(formatBytes))
    } else {
        s = fmt.Sprintf(string(formatBytes), a)
    }

    return s
}

func checkColor(color []byte) bool {
    check := false

    for _, validColor := range Colors {
        if strings.ToLower(string(color)) == strings.ToLower(string(validColor)) {
            check = true
        }
    }

    return check
}
