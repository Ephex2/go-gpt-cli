package image

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config/profile"
)

type imageEndpoint struct{}

var iEndpoint imageEndpoint

func (ce imageEndpoint) Name() string {
	return "image"
}

func (ie imageEndpoint) DefaultProfile() profile.Profile {
	p := ImageProfile{
		ProfileName:           "default",
		CreateImageBody:       GetDefaultCreateImageBody(),
		CreateDalle3ImageBody: GetDefaultDalle3Body(),
		CreateEditBody:        CreateImageEditBody,
		CreateVariationBody:   CreateVariationBody,
	}

	return p
}

func (ie imageEndpoint) ProfileFromJsonBuf(buf []byte) (p profile.Profile, err error) {
	var iProf ImageProfile
	err = json.Unmarshal(buf, &iProf)
	// if err != nil { return }

	p = iProf
	return
}

// This function is meant to be called by the init function in the audio cobra package
func Init() {
	profile.EndpointRegistry.Add(iEndpoint)
}
