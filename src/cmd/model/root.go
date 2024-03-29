package model

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ephex2/go-gpt-cli/log"
	"github.com/ephex2/go-gpt-cli/model"
	"github.com/spf13/cobra"
)

var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Allows you to interact with models hosted in the OpenAI API",
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Used to list all models",
	Long:    "Used to list all models hosted in your organization within theOpenAI API",
	Run:     listFunc,
	Args:    cobra.ExactArgs(0),
	Example: "go-gpt-cli model list",
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Retrieves a specific model",
	Long:    "Retrieves a specific model from within your organization with the OpenAI API",
	Run:     getFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli model get <modelName>",
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Deletes a specific model",
	Long:    "Deletes a specific model from within your organization with the OpenAI API",
	Run:     deleteFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli model delete <modelName>",
}

func listFunc(cmd *cobra.Command, args []string) {
	m, err := model.ListModels()
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func getFunc(cmd *cobra.Command, args []string) {
	m, err := model.RetrieveModel(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func deleteFunc(cmd *cobra.Command, args []string) {
	m, err := model.DeleteModel(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(m, "", "    ")
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
	//if err != nil { return nil }

	return
}

func init() {
	ModelCmd.AddCommand(listCmd)
	ModelCmd.AddCommand(getCmd)
	ModelCmd.AddCommand(deleteCmd)
}
