package audio

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ephex2/go-gpt-cli/audio"
	"github.com/ephex2/go-gpt-cli/log"
	"github.com/spf13/cobra"
)

var AudioCmd = &cobra.Command{
	Use:   "audio",
	Short: fmt.Sprintf("Allows you to make calls to the %s endpoint", audio.BaseRoute),
}

var speechCmd = &cobra.Command{
	Use:     "speech",
	Short:   "Used to create an audio file from text.",
	Run:     speechFunc,
	Args:    cobra.MinimumNArgs(1),
	Example: "go-gpt-cli audio speech \"Using quotes here is recommended. Inclusion of files is easy with $(cat filename)\"",
}

var promptCmd = &cobra.Command{
	Use:     "prompt",
	Short:   "Used get create a chat completion from a prompt, then read it over over the speaker. *Uses the default chat profile to generate the chat completion*",
	Run:     promptFunc,
	Args:    cobra.MinimumNArgs(1),
	Example: "go-gpt-cli audio prompt Feel free to add as many strings as you like but 'terminals will act best if you enclose your prompt in quotes'",
}

var playCmd = &cobra.Command{
	Use:     "play",
	Short:   "Plays an audio file from a provided path.",
	Long:    "Plays an audio file from a provided path. Only supports .mp3, .wav, and .flac files",
	Run:     playFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli audio play ./mySong.mp3",
}

var transcriptionCmd = &cobra.Command{
	Use:     "transcript",
	Short:   "Sends an audio file to the API in order to generate a text transcript. A prompt can be included in order to specify the textual 'style' of the transcript.",
	Run:     transcriptFunc,
	Args:    cobra.MinimumNArgs(1),
	Example: "go-gpt-cli audio transcript ./mySong.mp3 'This will be the text 'style' for the transcript. It should include proper punctuation, capitalization, and include the use of commas when appropriate'",
}

var verboseTranscriptionCmd = &cobra.Command{
	Use:     "verbose_transcript",
	Short:   "Same as 'transcript', but returns a much more verbose API response, including details about the segments of the transcribed text.",
	Run:     verboseTranscriptFunc,
	Args:    cobra.MinimumNArgs(1),
	Example: "go-gpt-cli audio transcript ./mySong.mp3 'This will be the text 'style' for the transcript. It should include proper punctuation, capitalization, and include the use of commas when appropriate'",
}

var translationCmd = &cobra.Command{
	Use:     "translation",
	Short:   "Sends an audio file to the API in order to generate a text translation to English. A prompt can be included in order to specify the textual 'style' of the transcript.",
	Run:     translationFunc,
	Args:    cobra.MinimumNArgs(1),
	Example: "go-gpt-cli audio translation ./mySpeech.mp3 'This will be the text 'style' for the transcript. It should include proper punctuation, capitalization, and include the use of commas when appropriate'",
}

var verboseTranslationCmd = &cobra.Command{
	Use:     "verbose_translation",
	Short:   "Same as 'translation', but returns a much more verbose API response, including details about the segments of the transcribed text.",
	Run:     verboseTranslationFunc,
	Args:    cobra.MinimumNArgs(1),
	Example: "go-gpt-cli audio translation ./mySpeech.mp3 'This will be the text 'style' for the transcript. It should include proper punctuation, capitalization, and include the use of commas when appropriate'",
}

func Execute(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		cmd.Help()
	}

	err = cmd.Execute()

	return
}

func speechFunc(cmd *cobra.Command, args []string) {
	s, err := audio.CreateSpeech(args)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(s)
}

func promptFunc(cmd *cobra.Command, args []string) {
	err := audio.ReadAudioPrompt(args)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}
}

func playFunc(cmd *cobra.Command, args []string) {
	err := audio.PlayAudioFile(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}
}

func transcriptFunc(cmd *cobra.Command, args []string) {
	var err error
	var s string

	if len(args) > 1 {
		s, err = audio.CreateTranscription(args[0], args[1:])
	} else {
		s, err = audio.CreateTranscription(args[0], []string{})
	}

	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(s)
}

func verboseTranscriptFunc(cmd *cobra.Command, args []string) {
	var err error
	var resp audio.CreateVerboseTranscriptionResponse

	if len(args) > 1 {
		resp, err = audio.CreateVerboseTranscription(args[0], args[1:])
	} else {
		resp, err = audio.CreateVerboseTranscription(args[0], []string{})
	}

	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Warning("Unable to marshal json from response, here is the unmarshalled response:")
		fmt.Println(resp)
	}

	fmt.Println(string(buf))
}

func translationFunc(cmd *cobra.Command, args []string) {
	var err error
	var s string

	if len(args) > 1 {
		s, err = audio.CreateTranslation(args[0], args[1:])
	} else {
		s, err = audio.CreateTranslation(args[0], []string{})
	}

	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	fmt.Println(s)
}

func verboseTranslationFunc(cmd *cobra.Command, args []string) {
	var err error
	var resp audio.CreateVerboseTranslationResponse

	if len(args) > 1 {
		resp, err = audio.CreateVerboseTranslation(args[0], args[1:])
	} else {
		resp, err = audio.CreateVerboseTranslation(args[0], []string{})
	}

	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	buf, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Warning("Unable to marshal json from response, here is the unmarshalled response:")
		fmt.Println(resp)
	}

	fmt.Println(string(buf))
}

func init() {
	AudioCmd.AddCommand(promptCmd)
	AudioCmd.AddCommand(playCmd)
	AudioCmd.AddCommand(speechCmd)
	AudioCmd.AddCommand(transcriptionCmd)
	AudioCmd.AddCommand(verboseTranscriptionCmd)
	AudioCmd.AddCommand(translationCmd)
	AudioCmd.AddCommand(verboseTranslationCmd)

	audio.Init()
}
