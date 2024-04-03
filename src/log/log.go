package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ephex2/go-gpt-cli/color"
)

const (
    LevelCritical = iota
    LevelWarning
    LevelInfo
    LevelDebug
)

var logLevel = LevelCritical

var logWriter io.Writer

var debugLogger *log.Logger
var infoLogger *log.Logger
var warningLogger *log.Logger
var criticalLogger *log.Logger


//var Logger *log.Logger
func SetLogLevel(level int) (e error) {
    switch level {
    case LevelCritical:
            logLevel = LevelCritical
        case LevelWarning:
            logLevel = LevelWarning
        case LevelInfo:
            logLevel = LevelInfo
        case LevelDebug:
            logLevel = LevelDebug
        default:
            e = InvalidLogError(strconv.Itoa(level))
    }

    return
}


func Debug(s string, a ...any)  {
    var input string

    if logLevel >= LevelDebug {
        if logWriter == os.Stdout {
            input = color.ColorSprintf(color.Blue, s, a)
        } else {
            input = fmt.Sprintf(s, a)
        }

        debugLogger.Print(input)
    }
}

func Info(s string, a ...any) {
    var input string

    if logLevel >= LevelInfo {
        if logWriter == os.Stdout {
            input = color.ColorSprintf(color.Green, s, a)
        } else {
            input = fmt.Sprintf(s, a)
        }

        infoLogger.Print(input)
    }
}

func Warning(s string, a ...any) {
    var input string

    if logLevel >= LevelWarning {
        if logWriter == os.Stdout {
            input = color.ColorSprintf(color.Yellow, s, a)
        } else {
            input = fmt.Sprintf(s, a)
        }

        warningLogger.Print(input)
    }
}

func Critical(s string, a ...any) {
    var input string

    if logLevel >= LevelCritical {
        if logWriter == os.Stdout {
            input = color.ColorSprintf(color.BrightRed, s, a)
        } else {
            input = fmt.Sprintf(s, a)
        }
    }

    criticalLogger.Print(input)
}

func init() {
    logWriter = os.Stdout

    debugLogger = log.New(logWriter, "DEBUG: ", log.Ldate|log.Ltime)
    infoLogger = log.New(logWriter, "INFO: ", log.Ldate|log.Ltime)
    warningLogger = log.New(logWriter, "WARING: ", log.Ldate|log.Ltime)
    criticalLogger = log.New(logWriter, "ERROR: ", log.Ldate|log.Ltime)
    /*
    var err error
    logWriter, err = os.Create("placeholder")
    if err != nil {
        panic(err.Error())
    }
    */
}
