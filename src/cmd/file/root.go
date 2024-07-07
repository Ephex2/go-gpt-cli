package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ephex2/go-gpt-cli/file"
	"github.com/ephex2/go-gpt-cli/log"
	"github.com/spf13/cobra"
)

var FileCmd = &cobra.Command{
	Use:   "file",
	Short: fmt.Sprintf("Allows you to make calls to the %s endpoint.", file.BaseFileRoute),
}

var createCmd = &cobra.Command{
	Use:               "create",
	Short:             "Used to create a file for use in assistants or fine-tuning",
	Long:              "Used to create a file for use in assistants or fine-tuning. Individual files cannot exceed 512MB or 2 million tokens for Assistants. The fine-tuning API only supports '.jsonl' files. An organization can upload a maximum of 100GB of files.",
	Run:               createFunc,
	Args:              cobra.MinimumNArgs(2),
	ValidArgsFunction: validFileCreateArgsFunc,
	Example:           "go-gpt-cli file create fine-tune myFile.jsonl",
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Used to delete a specific file uploaded to the vendor.",
	Long:    "Used to delete a specific file uploaded to the vendor. Must specify the ID when performing the call, which is returned for each file with the list operation",
	Run:     deleteFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli file delete fileId-abc123",
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Used to get a specific file uploaded to the vendor.",
	Long:    "Used to get a specific file uploaded to the vendor. Must specify the ID when performing the call, which is returned for each file with the list operation",
	Run:     getFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli file get fileId-abc123",
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Used to list all files uploaded to the vendor.",
	Long:    "Used to list all files uploaded to the vendor. This will list all files associated with the key's organization.",
	Run:     listFunc,
	Args:    cobra.ExactArgs(0),
	Example: "go-gpt-cli file list",
}

var statCmd = &cobra.Command{
	Use:     "stat",
	Short:   "Used to get stats of a specific file uploaded to the vendor.",
	Long:    "Used to get stats of a specific file uploaded to the vendor. Must specify the ID when performing the call, which is returned for each file with the list operation",
	Run:     statFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli file stat fileId-abc123",
}

func Execute(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		cmd.Help()
	}

	err = cmd.Execute()
	return
}

func createFunc(cmd *cobra.Command, args []string) {
	file, err := file.CreateFile(args[0], args[1])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(file, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func deleteFunc(cmd *cobra.Command, args []string) {
	del, err := file.DeleteFile(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(del, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func getFunc(cmd *cobra.Command, args []string) {
	buf, err := file.GetFile(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	// prints out actual file contents
	fmt.Println(string(buf))
}

func listFunc(cmd *cobra.Command, args []string) {
	files, err := file.ListFiles()
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(files, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func statFunc(cmd *cobra.Command, args []string) {
	file, err := file.StatFile(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(file, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func validFileCreateArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		// Not dynamic
		purposes := []string{file.AllowedFilePurposes.Assistants, file.AllowedFilePurposes.FineTune, file.AllowedFilePurposes.Batch}

		return purposes, cobra.ShellCompDirectiveNoFileComp
	}

	if len(args) == 1 {
		// let shell search for files
		return nil, cobra.ShellCompDirectiveDefault
	}

	// "Else"
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	FileCmd.AddCommand(createCmd)
	FileCmd.AddCommand(deleteCmd)
	FileCmd.AddCommand(getCmd)
	FileCmd.AddCommand(listCmd)
	FileCmd.AddCommand(statCmd)

	file.Init()
}
