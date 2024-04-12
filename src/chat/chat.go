package chat

import (
	"encoding/json"
	"errors"

	"github.com/ephex2/go-gpt-cli/api"
	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/image"
	"github.com/ephex2/go-gpt-cli/log"
)

func CreateChatCompletion(prompt []string) (content string, err error) {
	// Take user input and return completion completionConfig for request
	fPrompt := formatChat(prompt)
	log.Debug("Formatted chat string is : %s\n", fPrompt)

	chatProfile, err := getDefaultProfileFromPrompt(fPrompt)
	if err != nil {
		return
	}

	bufConfig, err := json.Marshal(chatProfile.CreateCompletionBody)
	if err != nil {
		return
	}

	log.Debug("Config is : %s\n", string(bufConfig))

	buf, err := api.GenericRequest(nil, bufConfig, "/v1/chat/completions", "POST", chatProfile.OverrideUrl())
	if err != nil {
		return
	}

	// Return completion response object.
	var completionResponse CompletionResponse
	err = json.Unmarshal(buf, &completionResponse)
	if err != nil {
		err = errors.New("unable to parse completion response.\nError is: " + err.Error())
		return
	} else if len(completionResponse.Choices) < 1 {
		err = errors.New("no choices returned for completion prompt")
		return
	}

	content = completionResponse.Choices[0].Message.Content
	err = postProcessing(completionResponse, chatProfile, "completion")
	if err != nil {
		return
	}

	return
}

// At time of creation, Open AI API supports png, jpg / jpeg, webp, and gif images
func CreateVisionChatCompletion(imagePath string, prompt []string) (resp string, err error) {
	b64Image, err := image.GetB64Encoding(imagePath)
	if err != nil {
		return
	}

	msg := formatChat(prompt)

	p, err := getDefaultVisionProfileFromPrompt(b64Image, msg)
	if err != nil {
		return
	}

	bufConfig, err := json.Marshal(p.CreateVisionCompletionBody)
	if err != nil {
		return
	}

	log.Debug("Config is : %s\n", string(bufConfig))

	buf, err := api.GenericRequest(nil, bufConfig, "/v1/chat/completions", "POST", p.OverrideUrl())
	if err != nil {
		return
	}

	var completionResponse CompletionResponse
	err = json.Unmarshal(buf, &completionResponse)
	if err != nil {
		err = errors.New("unable to parse completion response.\nError is: " + err.Error())
		return
	} else if len(completionResponse.Choices) < 1 {
		err = errors.New("no choices returned for completion prompt")
		return
	}

	resp = completionResponse.Choices[0].Message.Content
	err = postProcessing(completionResponse, p, "vision")
	if err != nil {
		return
	}

	return
}

// Used to clear all historical messages when message history is enabled in the profile. Should have no effect when the chat profile does not support history,
func ClearMessageHistory() (err error) {
    chatProfile := ChatProfile{}

	defaultProfileName, err := config.RuntimeConfig.GetDefaultProfile(chatProfile.Endpoint().Name())
	if err != nil {
		return
	}

	err = chatProfile.Load(defaultProfileName)
	if err != nil {
		return
	}

    err = chatProfile.ClearMessageHistory()
	if err != nil {
		return
	}

    return
}


func formatChat(chat []string) string {
	var formattedChat string

	for i, word := range chat {
		if i+1 == len(chat) {
			formattedChat += word
		} else {
			formattedChat += word + " "
		}
	}

	return formattedChat
}

func postProcessing(res CompletionResponse, profile ChatProfile, messageType string) (err error) {
	if profile.MessageHistory {
		if len(res.Choices) > 1 {
			err = errors.New("adding messages to history for responses with > 1 choice is not implemented")
			return
		} else if len(res.Choices) == 1 {
			switch messageType {
			case "completion":
				profile.AddCompletionMessage(res.Choices[0].Message)
			case "vision":
				msg := VisionMessage{
					Role: "assistant",
					Content: []VisionContent{
						{
							Type: "text",
							Text: &res.Choices[0].Message.Content,
						},
					},
				}

				profile.AddVisionMessage(msg)

			default:
				err = errors.New("message type not supported for post-processing. Message type is: " + messageType)
				return
			}

		} else {
			err = errors.New("no choices returned in response's CreateCompletionBody for post-processing")
			return
		}
	}

	return
}
