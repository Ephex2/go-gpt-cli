package finetuning

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/config/profile"
)

type FineTuningProfile struct {
	ProfileName        string
	CreateFineTuneBody CreateFineTuneBody
    Url                string
}

func (p FineTuningProfile) Name() string {
	if p.ProfileName == "" {
		p.Endpoint().DefaultProfile().Name()
	}

	return p.ProfileName
}

func (p FineTuningProfile) OverrideUrl() string {
    return p.Url
}

func (p FineTuningProfile) SetName(name string) profile.Profile {
	p.ProfileName = name
	return p
}

func (p FineTuningProfile) ProfileRepository() profile.Repository {
	return profile.RuntimeRepository
}

func (p FineTuningProfile) Endpoint() profile.Endpoint {
	return fEndpoint
}

func (p *FineTuningProfile) Load(profileName string) (err error) {
	buf, err := profile.Repository.Read(p.ProfileRepository(), profileName, p.Endpoint().Name())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, p)
	//if err != nil { return }

	return
}

// Use this function to generate the initial config for a profile
func getDefaultProfile() (profile FineTuningProfile, err error) {
	defaultProfileName, err := config.RuntimeConfig.GetDefaultProfile(FineTuningProfile{}.Endpoint().Name())
	if err != nil {
		return
	}

	err = profile.Load(defaultProfileName)
	if err != nil {
		return
	}

	return
}
