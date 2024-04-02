package chat

import (
	"fmt"
	"os"

	"github.com/ephex2/go-gpt-cli/chat"
	"github.com/ephex2/go-gpt-cli/log"
	"github.com/spf13/cobra"
)

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Allows you to make calls to the /v1/chat/ endpoint",
}


var promptCmd = &cobra.Command{
	Use:     "prompt",
	Short:   "Used to get a prompt from a chat completion, then create a .mp3 file using the audio create speech endpoint. *Reads the file over the speaker*",
	Run:     promptFunc,
	Args:    cobra.MinimumNArgs(1),
	Example: "go-gpt-cli chat prompt Feel free to add as many strings as you like but 'terminals will act best if you enclose your prompt in quotes'",
}


var visionCmd = &cobra.Command{
	Use:     "vision",
	Short:   "Used to make vision requests to multi-modal LLMs.",
	Long:    "Used to make vision requests to multi-modal LLMs. The first argument is a path to a valid image file, all other arguments are concatenated as a prompt.",
	Run:     visionFunc,
	Args:    cobra.MinimumNArgs(2),
	Example: "go-gpt-cli chat vision ./myimage.png'terminals will act best if you enclose your prompt in quotes'",
}


var clearCmd = &cobra.Command{
    Use: "clear",
	Short:   "Used to clear all historical messages when message history is enabled in the profile.",
	Long:   "Used to clear all historical messages when message history is enabled in the profile. Should have no effect when the chat profile does not support history",
	Run:     clearFunc,
	Args:    cobra.ExactArgs(0),
	Example: "go-gpt-cli chat clear",
    
}


func Execute(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		cmd.Help()
	}

	err = cmd.Execute()
	return
}


func promptFunc(cmd *cobra.Command, args []string) {
	s, err := chat.CreateChatCompletion(args)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(s)
}


func visionFunc(cmd *cobra.Command, args []string) {
	s, err := chat.CreateVisionChatCompletion(args[0], args[1:])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(s)
}


func clearFunc(cmd *cobra.Command, args []string) {
    err := chat.ClearMessageHistory()
    if err != nil {
        log.Critical(err.Error() + "\n")
        os.Exit(1)
    }
}


func init() {
    ChatCmd.AddCommand(clearCmd)
    ChatCmd.AddCommand(promptCmd)
    ChatCmd.AddCommand(visionCmd)
    chat.Init()
}
