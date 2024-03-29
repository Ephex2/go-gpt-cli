package embeddings

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/config/profile"
)

type EmbeddingsProfile struct {
	ProfileName         string
	CreateEmbeddingBody CreateEmbeddingBody
}

func (p EmbeddingsProfile) Name() string {
	if p.ProfileName == "" {
		p.Endpoint().DefaultProfile().Name()
	}

	return p.ProfileName
}

func (p EmbeddingsProfile) SetName(name string) profile.Profile {
	p.ProfileName = name
	return p
}

func (p EmbeddingsProfile) ProfileRepository() profile.Repository {
	return profile.RuntimeRepository
}

func (p EmbeddingsProfile) Endpoint() profile.Endpoint {
	return eEndpoint
}

func (p *EmbeddingsProfile) Load(profileName string) (err error) {
	buf, err := profile.Repository.Read(p.ProfileRepository(), profileName, p.Endpoint().Name())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, p)
	//if err != nil { return }

	return
}

// Use this function to generate the initial config for a profile
func getDefaultProfile() (profile EmbeddingsProfile, err error) {
	defaultProfileName, err := config.RuntimeConfig.GetDefaultProfile(EmbeddingsProfile{}.Endpoint().Name())
	if err != nil {
		return
	}

	err = profile.Load(defaultProfileName)
	if err != nil {
		return
	}

	return
}
