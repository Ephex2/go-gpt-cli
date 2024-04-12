package audio

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config/profile"
)

const EndpointName = "audio"

type audioEndpoint struct{}

var aEndpoint audioEndpoint

func (ae audioEndpoint) Name() string {
	return "audio"
}

func (ae audioEndpoint) DefaultProfile() profile.Profile {
	p := AudioProfile{
		ProfileName:                    "default",
		CreateSpeechBody:               DefaultCreateSpeechBody(),
		CreateTranscriptionBody:        CreateTranscriptionBody,
		CreateVerboseTranscriptionBody: CreateVerboseTranscriptionBody,
		CreateTranslationBody:          CreateTranslationBody,
		CreateVerboseTranslationBody:   CreateVerboseTranslationBody,
		SaveDirectory:                  "",
        Url: "",
	}

	return p
}

func (ae audioEndpoint) ProfileFromJsonBuf(buf []byte) (ap profile.Profile, err error) {
	var audioP AudioProfile
	err = json.Unmarshal(buf, &audioP)
	if err != nil {
		return
	}

	ap = audioP
	return
}

// Why Init() ?
// The initialization of main seemed to only call init() in packages imported by main (or by packages imported by packages imported by main).
// This function is meant to be called by the init function in the audio cobra package
// This may be due to 'skill-issues' -- consider renaming to init() and removing calls to audio.Init() if this is the case.
func Init() {
	profile.EndpointRegistry.Add(aEndpoint)
}
