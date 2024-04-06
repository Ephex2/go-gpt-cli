package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/log"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Allows you to interact with global settings",
}

var setKeyCmd = &cobra.Command{
	Use:     "apikey",
	Short:   "Used to set the apikey which will be used when authenticating to API endpoints. It is not encrypted or otherwise protected at rest.",
	Run:     setKeyFunc,
	Args:    cobra.ExactArgs(1),
    Aliases: []string{"setkey"},
	Example: "go-gpt-cli config apikey 12345",
}

var setUrlCmd = &cobra.Command{
	Use:     "seturl",
	Short:   "Used to set the base url which will be used when calling API endpoints. Defaults to the OpenAI API url.",
	Run:     setUrlFunc,
	Args:    cobra.ExactArgs(1),
    Example: "go-gpt-cli config seturl http://my.alternative.name:port",
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Lists the current global settings",
	Run:     getFunc,
	Args:    cobra.ExactArgs(0),
	Example: "go-gpt-cli config get",
}

func setKeyFunc(cmd *cobra.Command, args []string) {
	err := config.SetApiKey(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}
}

func setUrlFunc(cmd *cobra.Command, args []string) {
	err := config.SetBaseUrl(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}
}

func getFunc(cmd *cobra.Command, args []string) {
	tempConfig := config.RuntimeConfig

	buf, err := json.MarshalIndent(tempConfig.Settings, "", "  ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func Execute(cmd *cobra.Command, args []string) (err error) {

	if len(args) < 1 {
		cmd.Help()
	}

	err = cmd.Execute()

	return
}

func init() {
	ConfigCmd.AddCommand(setKeyCmd)
	ConfigCmd.AddCommand(setUrlCmd)
	ConfigCmd.AddCommand(getCmd)
}
