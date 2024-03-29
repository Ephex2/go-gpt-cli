package main

import (
"github.com/ephex2/go-gpt-cli/cmd"
"github.com/ephex2/go-gpt-cli/log"
)

func main() {
    err := cmd.Execute()
    if err != nil {
        log.Critical(err.Error() + "\n")
    }
}
