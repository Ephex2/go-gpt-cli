package log

import (
	"github.com/ephex2/go-gpt-cli/color"
)

const (
    critical = iota
    warning
    info
    debug
)

var logLevel = info

func SetLogLevel(s string) (e error) {
    switch s {
        case "debug":
            logLevel = debug
        case "info":
            logLevel = info
        case "warning":
            logLevel = warning
        case "critical":
            logLevel = critical
        default:
            e = InvalidLogError(s)
    }

    return
}

func Debug(s string, a ...any)  {
    if logLevel >= debug {
        color.ColorPrintf(color.Blue, s, a...)
    }
}

func Info(s string, a ...any) {
    if logLevel >= info {
       color.ColorPrintf(color.Green, s, a...)
    }
}

func Warning(s string, a ...any) {
    if logLevel >= warning {
       color.ColorPrintf(color.Yellow, s, a...)
    }
}

func Critical(s string, a ...any) {
    if logLevel >= critical {
       color.ColorPrintf(color.Red, s, a...)
    }
}
