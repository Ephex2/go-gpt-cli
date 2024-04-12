package file

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config/profile"
)

type fileEndpoint struct{}

var fEndpoint fileEndpoint

func (e fileEndpoint) Name() string {
	return "file"
}

// This is to be used when creating profiles for the first time.
func (e fileEndpoint) DefaultProfile() profile.Profile {
	p := FileProfile{
		ProfileName:    "default",
		CreateFileBody: GetDefaultBody(),
        Url: "",
	}

	return p
}

func (e fileEndpoint) ProfileFromJsonBuf(buf []byte) (p profile.Profile, err error) {
	var cProf FileProfile
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
