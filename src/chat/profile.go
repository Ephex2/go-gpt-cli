package chat

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/config/profile"
)

type ChatProfile struct {
	ProfileName          string
	CreateCompletionBody       CreateCompletionBody
	CreateVisionCompletionBody CreateVisionCompletionBody
	MessageHistory       bool
}

func (c ChatProfile) Name() string {
	if c.ProfileName == "" {
		return "default"
	}

	return c.ProfileName
}

func (c ChatProfile) SetName(name string) profile.Profile {
	c.ProfileName = name
	return c
}

func (c ChatProfile) ProfileRepository() profile.Repository {
	return profile.RuntimeRepository
}

func (c ChatProfile) Endpoint() profile.Endpoint {
	return cEndpoint
}

func (c *ChatProfile) AddCompletionMessage(msg Message) (err error) {
	c.CreateCompletionBody.Messages = append(c.CreateCompletionBody.Messages, msg)

	if c.MessageHistory {
		err = c.ProfileRepository().Update(c)
	}

	return
}

func (c *ChatProfile) AddVisionMessage(msg VisionMessage) (err error) {
	c.CreateVisionCompletionBody.Messages = append(c.CreateVisionCompletionBody.Messages, msg)

	if c.MessageHistory {
		err = c.ProfileRepository().Update(c)
	}

	return
}

func (c *ChatProfile) Load(profileName string) (err error) {
	buf, err := c.ProfileRepository().Read(profileName, c.Endpoint().Name())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, c)
	//if err != nil { return }

	return
}

// Keep all system messages by defauly, assuming that they are customized to initialize prompts / conversations.
func (c *ChatProfile) ClearMessageHistory() (err error) {
	err = c.Load(c.Name())
	if err != nil {
		return
	}

	var defaultSystemMessage []Message
	for _, message := range c.CreateCompletionBody.Messages {
		if strings.ToLower(message.Role) == "system" {
			defaultSystemMessage = append(defaultSystemMessage, message)
		}
	}

	if len(defaultSystemMessage) == len(c.CreateCompletionBody.Messages) {
		// nothing to do, return
		return
	}

	c.CreateCompletionBody.Messages = defaultSystemMessage
	err = c.ProfileRepository().Update(c)
	return
}

// Use this function to generate the initial config for a profile when taking string input for 'normal' chat completions
func getDefaultProfileFromPrompt(prompt string) (profile ChatProfile, err error) {
	if prompt == "" {
		err = errors.New("please provide a prompt when building a profile from prompt")
		return
	}

	defaultProfileName, err := config.RuntimeConfig.GetDefaultProfile(ChatProfile{}.Endpoint().Name())
	if err != nil {
		return
	}

	err = profile.Load(defaultProfileName)
	if err != nil {
		return
	}

	msg := Message{Role: "user", Content: prompt}
	err = profile.AddCompletionMessage(msg)
	if err != nil {
		return
	}

	return
}

// Use this function to generate the initial config for a profile when taking string input
func getDefaultVisionProfileFromPrompt(b64image string, prompt string) (profile ChatProfile, err error) {
	if prompt == "" {
		err = errors.New("please provide a prompt when building a profile from prompt")
		return
	}

	defaultProfileName, err := config.RuntimeConfig.GetDefaultProfile(ChatProfile{}.Endpoint().Name())
	if err != nil {
		return
	}

	err = profile.Load(defaultProfileName)
	if err != nil {
		return
	}

	var content []VisionContent
	if b64image == "" {
		content = []VisionContent{
			{
				Type: "text",
				Text: &prompt,
			},
		}
	} else {
		content = []VisionContent{
			{
				Type: "text",
				Text: &prompt,
			},
			{
				Type: "image_url",
				ImageUrl: ImageUrl{
					Url: b64image,
				},
			},
		}
	}

	msg := VisionMessage{Role: "user", Content: content}
	err = profile.AddVisionMessage(msg)
	if err != nil {
		return
	}

	return
}
