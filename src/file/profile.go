package file

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/config/profile"
)

type FileProfile struct {
	ProfileName    string
	CreateFileBody map[string]string
    Url            string
}

func (p FileProfile) Name() string {
	if p.ProfileName == "" {
		p.Endpoint().DefaultProfile().Name()
	}

	return p.ProfileName
}

func (p FileProfile) OverrideUrl() string {
    return p.Url
}

func (p FileProfile) SetName(name string) profile.Profile {
	p.ProfileName = name
	return p
}

func (p FileProfile) ProfileRepository() profile.Repository {
	return profile.RuntimeRepository
}

func (p FileProfile) Endpoint() profile.Endpoint {
	return fEndpoint
}

func (p *FileProfile) Load(profileName string) (err error) {
	buf, err := profile.Repository.Read(p.ProfileRepository(), profileName, p.Endpoint().Name())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, p)
	//if err != nil { return }

	return
}

// Use this function to generate the initial config for a profile
func getDefaultProfile() (profile FileProfile, err error) {
	defaultProfileName, err := config.RuntimeConfig.GetDefaultProfile(FileProfile{}.Endpoint().Name())
	if err != nil {
		return
	}

	err = profile.Load(defaultProfileName)
	if err != nil {
		return
	}

	return
}
