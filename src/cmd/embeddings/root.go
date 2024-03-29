package embeddings

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ephex2/go-gpt-cli/embeddings"
	"github.com/ephex2/go-gpt-cli/log"
	"github.com/spf13/cobra"
)

var EmbeddingsCmd = &cobra.Command{
	Use:   "embeddings",
	Short: fmt.Sprintf("Allows you to make calls to the %s endpoint.", embeddings.BaseEmbeddingsRoute),
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Used to create embeddings from an array of strings input.",
	Long:    "Used to create embeddings from an array of strings input. Each argument passed to the CLI separated by a space will be a separate element in the array.",
	Run:     createFunc,
	Args:    cobra.MinimumNArgs(1),
	Example: "go-gpt-cli embeddings create Feel free to add as many strings as you like but 'put strings in quotes if you want them to be one element in the array'",
}

func Execute(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		cmd.Help()
	}

	err = cmd.Execute()
	return
}

func createFunc(cmd *cobra.Command, args []string) {
	res, err := embeddings.CreateEmbeddings(args)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	// Outputting embedding objects as string for now
	for _, embedding := range res.Data {
		buf, err := json.MarshalIndent(embedding, "", "    ")
		if err != nil {
			log.Critical(err.Error() + "\n")
			os.Exit(1)
		}

		fmt.Println(string(buf))
	}
}

func init() {
	EmbeddingsCmd.AddCommand(createCmd)

	embeddings.Init()
}
