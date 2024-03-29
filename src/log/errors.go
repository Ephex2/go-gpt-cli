package log

import (
    "fmt"
)


type InvalidLogError string

func (e InvalidLogError) Error() string {
    return fmt.Sprintf("log: Invalid LogLevel: %s. Please use one of [Debug,Info,Warning,Critical]", e)
}

func (e InvalidLogError) Timeout() bool {
    return false
}

func (e InvalidLogError) Temporary() bool {
    return false
}

