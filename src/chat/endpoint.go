package chat

import (
	//	"github.com/ephex2/go-gpt-cli/config/repository"
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config/profile"
)

type chatEndpoint struct{}

var cEndpoint chatEndpoint

func (ce chatEndpoint) Name() string {
	return "chat"
}

// Define the default profile configuration for the chat endpoint
// This is to be used when creating profiles for the first time.
func (ce chatEndpoint) DefaultProfile() profile.Profile {
	p := ChatProfile{
		ProfileName:          "default",
		CreateCompletionBody:       GetDefaultBody(),
		CreateVisionCompletionBody: GetDefaultVisionBody(),
		MessageHistory:       false,
        Url: "",
	}

	return p
}

func (ce chatEndpoint) ProfileFromJsonBuf(buf []byte) (p profile.Profile, err error) {
	var cProf ChatProfile
	err = json.Unmarshal(buf, &cProf)
	if err != nil {
		return
	}

	p = cProf
	return
}

// Why Init() ? The initialization of main seemed to only call init in packages imported by main.
// This function is meant to be called by the init function in the audio cobra package
func Init() {
	profile.EndpointRegistry.Add(cEndpoint)
}
