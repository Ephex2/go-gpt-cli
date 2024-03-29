package finetuning

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config/profile"
)

type finetuningEndpoint struct{}

var fEndpoint finetuningEndpoint

func (e finetuningEndpoint) Name() string {
	return "finetuning"
}

// This is to be used when creating profiles for the first time.
func (e finetuningEndpoint) DefaultProfile() profile.Profile {
	p := FineTuningProfile{
		ProfileName:           "default",
		CreateFineTuneRequest: DefaultCreateFineTuneRequest(),
	}

	return p
}

func (e finetuningEndpoint) ProfileFromJsonBuf(buf []byte) (p profile.Profile, err error) {
	var cProf FineTuningProfile
	err = json.Unmarshal(buf, &cProf)
	// if err != nil { return }

	p = cProf
	return
}

// Why Init() ? The initialization of main seemed to only call init in packages imported by main.
// This function is meant to be called by the init function in the audio cobra package
func Init() {
	profile.EndpointRegistry.Add(fEndpoint)
}
