package finetuning

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ephex2/go-gpt-cli/finetuning"
	"github.com/ephex2/go-gpt-cli/log"
	"github.com/spf13/cobra"
)

var FineTuningCmd = &cobra.Command{
	Use:   "finetuning",
	Short: fmt.Sprintf("Allows you to make calls to the %s endpoint.", finetuning.BaseFineTuningRoute),
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Used to create a finetuning job for a model",
	Long:    "Used to create a finetuning job for a model. If no training file is specified, uses the training file from the default profile.\nThe job will not complete immediately, and other commands can be used to check the status of the job and the output model.",
	Run:     createFunc,
	Args:    cobra.MaximumNArgs(1),
	Example: "go-gpt-cli finetuning create trainingFile",
}

var cancelCmd = &cobra.Command{
	Use:     "cancel",
	Short:   "Used to cancel a specific finetuning job.",
	Long:    "Used to cancel a specific finetuning job. Must specify the ID when performing the call, which is returned for each finetuning job with the jobs operation",
	Run:     cancelFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli finetuning cancel fileId-abc123",
}

var eventsCmd = &cobra.Command{
	Use:     "events",
	Short:   "Used to list all events for a specific job id.",
	Long:    "Used to list all events for a specific job id. Is limited to your organization when calling the Open AI api.",
	Run:     eventsFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli finetuning events ftjob-abc123",
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Used to get a specific finetuning job started by the vendor.",
	Long:    "Used to get a specific finetuning job started by the vendor. Must specify the ID when performing the call, which is returned for each finetuning job with the jobs command",
	Run:     getFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli finetuning get fileId-abc123",
}

var jobsCmd = &cobra.Command{
	Use:     "jobs",
	Short:   "Used to list all jobs running.",
	Long:    "Used to list all jobs running. Is limited to your organization when calling the Open AI api.",
	Run:     jobsFunc,
	Args:    cobra.ExactArgs(0),
	Example: "go-gpt-cli finetuning jobs",
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
	if len(args) == 1 {
		id = args[0]
	}

	finetuning, err := finetuning.CreateJob(id)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(finetuning, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func cancelFunc(cmd *cobra.Command, args []string) {
	job, err := finetuning.CancelJob(args[0])
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

func eventsFunc(cmd *cobra.Command, args []string) {
	events, err := finetuning.ListEvents(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(events, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(string(buf))
}

func getFunc(cmd *cobra.Command, args []string) {
	job, err := finetuning.GetJob(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(job, "", "    ")
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	// prints out actual finetuning contents
	fmt.Println(string(buf))
}

func jobsFunc(cmd *cobra.Command, args []string) {
	jobs, err := finetuning.ListJobs()
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
	FineTuningCmd.AddCommand(cancelCmd)
	FineTuningCmd.AddCommand(createCmd)
	FineTuningCmd.AddCommand(eventsCmd)
	FineTuningCmd.AddCommand(getCmd)
	FineTuningCmd.AddCommand(jobsCmd)

	finetuning.Init()
}
