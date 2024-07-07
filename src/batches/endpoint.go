package batches


import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config/profile"
)

type batchEndpoint struct{}

var bEndpoint batchEndpoint

func (e batchEndpoint) Name() string {
	return "batch"
}

// This is to be used when creating profiles for the first time.
func (e batchEndpoint) DefaultProfile() profile.Profile {
	p := BatchProfile{
		ProfileName:    "default",
		CreateBatchBody: GetDefaultBody(),
        Url: "",
	}

	return p
}

func (e batchEndpoint) ProfileFromJsonBuf(buf []byte) (p profile.Profile, err error) {
	var cProf BatchProfile
	err = json.Unmarshal(buf, &cProf)
	if err != nil { return }

	p = cProf
	return
}

// Why Init() ? The initialization of main seemed to only call init in packages imported by main.
// This function is meant to be called by the init function in the audio cobra package
func Init() {
	profile.EndpointRegistry.Add(bEndpoint)
}
