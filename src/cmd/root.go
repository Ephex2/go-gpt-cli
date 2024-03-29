package cmd

import (
	"github.com/ephex2/go-gpt-cli/cmd/audio"
	"github.com/ephex2/go-gpt-cli/cmd/chat"
	"github.com/ephex2/go-gpt-cli/cmd/config"
	"github.com/ephex2/go-gpt-cli/cmd/embeddings"
	"github.com/ephex2/go-gpt-cli/cmd/file"
	"github.com/ephex2/go-gpt-cli/cmd/finetuning"
	"github.com/ephex2/go-gpt-cli/cmd/image"
	"github.com/ephex2/go-gpt-cli/cmd/model"
	"github.com/ephex2/go-gpt-cli/cmd/profile"
	"github.com/ephex2/go-gpt-cli/config/repository"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-gpt-cli",
	Short: "A CLI tool for interacting with the open AI API",
}

func Execute() error {
	rootCmd.AddCommand(audio.AudioCmd)
	rootCmd.AddCommand(chat.ChatCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(embeddings.EmbeddingsCmd)
	rootCmd.AddCommand(file.FileCmd)
	rootCmd.AddCommand(finetuning.FineTuningCmd)
	rootCmd.AddCommand(image.ImageCmd)
	rootCmd.AddCommand(model.ModelCmd)
	rootCmd.AddCommand(profile.ProfileCmd)

	err := rootCmd.Execute()
	return err
}

func init() {
	repository.Init()
}
