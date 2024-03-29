package image

import (
	"fmt"
	"os"
	"sort"

	"github.com/ephex2/go-gpt-cli/image"
	"github.com/ephex2/go-gpt-cli/log"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
)

var ImageCmd = &cobra.Command{
	Use:   "image",
	Short: fmt.Sprintf("Allows you to make calls to the %s endpoint.", image.BaseImageRoute),
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Used to create images from a prompt. The first string provided will be the filepath for output. Other strings are the prompt.",
	Long:    "Used to create images from an array of strings input. The first string is the filepath where the image will be written. Note that since multiple images can be generated per request to the API, this may return an array of paths.",
	Run:     createFunc,
	Args:    cobra.MinimumNArgs(2),
	Example: "go-gpt-cli image create myimage.png Feel free to add as many strings as you like but 'put strings in quotes if you want them to be one element in the array'",
}

var dalle3CreateCmd = &cobra.Command{
	Use:     "dalle3create",
	Short:   "Used to create images from a prompt using the dalle3 OpenAI model.",
	Long:    "Used to create images from an array of strings input using the dalle3 OpenAI model. The first string is the filepath where the image will be written. Note that since multiple images can be generated per request to the API, this may return an array of paths.",
	Run:     dalle3CreateFunc,
	Args:    cobra.MinimumNArgs(2),
	Example: "go-gpt-cli image create myimage.png Feel free to add as many strings as you like but 'put strings in quotes if you want them to be one element in the array'",
}

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Used to edit images from a given square image.",
	Long:    "Used to edit images from a given square image. A mask can be provided whose transparent areas (alpha is zero) determine where the image should be edited. The image file must be square, less than 4MB, and square for the OpenAI API. Image size will be determined dynamically at runtime.",
	Run:     editFunc,
	Args:    cobra.MinimumNArgs(2),
	Example: "go-gpt-cli image edit myimage.png mymask.png 'A prompt describing the resulting edited images'",
}

var variationCmd = &cobra.Command{
	Use:     "variation",
	Short:   "Used to create image variations from a given square image.",
	Long:    "Used to create image variations from a given square image. The image file must be square, less than 4MB, and square for the OpenAI API. Image size will be determined dynamically at runtime.",
	Run:     variationFunc,
	Args:    cobra.ExactArgs(1),
	Example: "go-gpt-cli image variation myimage.png",
}

func Execute(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		cmd.Help()
	}

	err = cmd.Execute()
	return
}

func createFunc(cmd *cobra.Command, args []string) {
	paths, revisedPrompt, err := image.CreateImage(args[0], args[1:])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	if revisedPrompt != "" {
		fmt.Println("Prompt was revised to: " + revisedPrompt)
	}

	outputPaths(paths)
}

func dalle3CreateFunc(cmd *cobra.Command, args []string) {
	paths, revisedPrompt, err := image.CreateDalle3Image(args[0], args[1:])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	if revisedPrompt != "" {
		fmt.Println("Prompt was revised to: " + revisedPrompt)
	}

	outputPaths(paths)
}

func editFunc(cmd *cobra.Command, args []string) {
	var prompt []string
	var mask *os.File
	mask = nil

	// mask setup
	if maskStat, err := os.Stat(args[1]); err == nil {
		if !maskStat.IsDir() {
			t, err := mimetype.DetectFile(args[1])
			if err == nil {
				if t.Extension() == ".png" {
					mask, err = os.Open(args[1])
					if err != nil {
						mask = nil
					}
				}
			}
		}
	}

	// prompt setup
	if mask != nil {
		if len(args) == 2 {
			log.Critical("Please provide a prompt for your image edit; only a base image and mask were found as arguments.\n")
			os.Exit(1)
		}

		prompt = append(prompt, args[2:]...)
	} else {
		prompt = append(prompt, args[1:]...)
	}

	paths, err := image.CreateEdit(args[0], mask, prompt)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	outputPaths(paths)
}

func variationFunc(cmd *cobra.Command, args []string) {
	paths, err := image.CreateVariation(args[0])
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	outputPaths(paths)
}

func outputPaths(paths []string) {
	sort.Strings(paths)

	log.Debug("Created images: ")

	for _, path := range paths {
		fmt.Println(path)
	}
}

func init() {
	ImageCmd.AddCommand(createCmd)
	ImageCmd.AddCommand(dalle3CreateCmd)
	ImageCmd.AddCommand(editCmd)
	ImageCmd.AddCommand(variationCmd)

	image.Init()
}
