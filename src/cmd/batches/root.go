package batches

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ephex2/go-gpt-cli/batches"
	"github.com/ephex2/go-gpt-cli/log"
	"github.com/spf13/cobra"
)

var BatchesCmd = &cobra.Command{
	Use:   "batches",
	Short: fmt.Sprintf("Allows you to make calls to the %s endpoint.", batches.BaseBatchesRoute),
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Used to create a batches job for a model",
	Long:    "Used to create a batches job for a given api endpoint. Must specify both a file containing jsonl data that has been uploaded, and an endpoint string.\nThe batch will not complete immediately, and other commands can be used to check the status of the batch.",
	Run:     createFunc,
	Args:    cobra.ExactArgs(2),
	Example: "go-gpt-cli batches create fileId targetEndPoint",
}

var cancelCmd = &cobra.Command{
	Use:     "cancel",
	Short:   "Used to cancel a specific batches job.",
	Long:    "Used to cancel a specific batches job. Must specify the ID when performing the call.",
	Run:     cancelFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli batches cancel batchId-abc123",
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Used to get a specific batches job started by the vendor.",
	Long:    "Used to get a specific batches job started by the vendor. Must specify the ID when performing the call, which is returned for each batches job with the jobs command",
	Run:     getFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli batches get batchId-abc123",
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Used to list all jobs running.",
	Long:    "Used to list all jobs running. Is limited to your organization when calling the OpenAI api.",
	Run:     listFunc,
	Args:    cobra.ExactArgs(0),
	Example: "go-gpt-cli batches list",
}

func Execute(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		cmd.Help()
	}

	err = cmd.Execute()
	return
}

func createFunc(cmd *cobra.Command, args []string) {
	var id string
    var endpoint string 

	id = args[0]
    endpoint = args[1]

	batches, err := batches.CreateBatch(id, endpoint)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(batches, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func cancelFunc(cmd *cobra.Command, args []string) {
	job, err := batches.CancelBatch(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(job, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func getFunc(cmd *cobra.Command, args []string) {
	job, err := batches.GetBatch(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(job, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	// prints out actual batches contents
	fmt.Println(string(buf))
}

func listFunc(cmd *cobra.Command, args []string) {
	jobs, err := batches.ListBatches()
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(jobs, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func init() {
	BatchesCmd.AddCommand(cancelCmd)
	BatchesCmd.AddCommand(createCmd)
	BatchesCmd.AddCommand(getCmd)
	BatchesCmd.AddCommand(listCmd)

	batches.Init()
}
