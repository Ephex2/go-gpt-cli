package embeddings

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config/profile"
)

type embeddingsEndpoint struct{}

var eEndpoint embeddingsEndpoint

func (e embeddingsEndpoint) Name() string {
	return "embeddings"
}

// This is to be used when creating profiles for the first time.
func (e embeddingsEndpoint) DefaultProfile() profile.Profile {
	p := EmbeddingsProfile{
		ProfileName:         "default",
		CreateEmbeddingBody: GetDefaultBody(),
        Url: "",
	}

	return p
}

func (e embeddingsEndpoint) ProfileFromJsonBuf(buf []byte) (p profile.Profile, err error) {
	var cProf EmbeddingsProfile
	err = json.Unmarshal(buf, &cProf)
	// if err != nil { return }

	p = cProf
	return
}

// Why Init() ? The initialization of main seemed to only call init in packages imported by main.
// This function is meant to be called by the init function in the audio cobra package
func Init() {
	profile.EndpointRegistry.Add(eEndpoint)
}
