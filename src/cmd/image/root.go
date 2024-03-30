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
	Short:   "Used to create images from a prompt.",
	Long:    "Used to create images from an array of strings input. The first string is the folder path where the image will be written. Note that since multiple images can be generated per request to the API, this may return an array of paths.",
	Run:     createFunc,
	Args:    cobra.MinimumNArgs(2),
	Example: "go-gpt-cli image create outputFolderPath 'image creation text prompt'",
}

var dalle3CreateCmd = &cobra.Command{
	Use:     "dalle3create",
	Short:   "Used to create images from a prompt using the dalle3 OpenAI model.",
	Long:    "Used to create images from an array of strings input using the dalle3 OpenAI model. The first string is the folder path where the image will be written. Note that since multiple images can be generated per request to the API, this may return an array of paths.",
	Run:     dalle3CreateFunc,
	Args:    cobra.MinimumNArgs(2),
	Example: "go-gpt-cli image dalle3create outputFolderPath 'image creation text prompt'",
}

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Used to edit images from a given square image.",
	Long:    "Used to edit images from a given square image. The first argument is the output folder path, the next is the path to the image to edit, the third is an optional path to a mask file, and the rest will be included a s aprompt to the API. A mask can be provided whose transparent areas (alpha is zero) determine where the image should be edited. The image file must be square, less than 4MB, and square for the OpenAI API. Image size will be determined dynamically at runtime.",
	Run:     editFunc,
	Args:    cobra.MinimumNArgs(3),
	Example: "go-gpt-cli image edit outputFolderPath myimage.png mymask.png 'A prompt describing the resulting edited images'",
}

var variationCmd = &cobra.Command{
	Use:     "variation",
	Short:   "Used to create image variations from a given square image.",
	Long:    "Used to create image variations from a given square image. The first argument is the output folder path, and the second one will be the image to create a variation for. The image file must be square, less than 4MB, and square for the OpenAI API. Image size will be determined dynamically at runtime.",
	Run:     variationFunc,
	Args:    cobra.ExactArgs(2),
	Example: "go-gpt-cli image variation outputFolderPath myimage.png",
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

    folderPath := args[0]
    imagePath := args[1]
    potentialMask := args[2] // third argument

	if maskStat, err := os.Stat(potentialMask); err == nil {
		if !maskStat.IsDir() {
			t, err := mimetype.DetectFile(potentialMask)
			if err == nil {
				if t.Extension() == ".png" {
					mask, err = os.Open(potentialMask)
					if err != nil {
						mask = nil
					}
				}
			}
		}
	}

	// prompt setup
	if mask != nil {
		if len(args) == 3 {
			log.Critical("Please provide a prompt for your image edit; only a base image and mask were found as arguments.\n")
			os.Exit(1)
		}

		prompt = append(prompt, args[3:]...)
	} else {
		prompt = append(prompt, args[2:]...)
	}

	paths, err := image.CreateEdit(imagePath, mask, folderPath, prompt)
	if err != nil {
		log.Critical(err.Error() + "\n")
		os.Exit(1)
	}

	outputPaths(paths)
}

func variationFunc(cmd *cobra.Command, args []string) {
    folderPath := args[0]
    filePath := args[1]

	paths, err := image.CreateVariation(filePath, folderPath)
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
