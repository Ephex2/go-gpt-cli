package image

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/config/profile"
)

type ImageProfile struct {
	ProfileName           string
	CreateImageBody       CreateImageBody
	CreateDalle3ImageBody CreateDalle3ImageBody
	CreateEditBody        map[string]string
	CreateVariationBody   map[string]string
}

func (ip ImageProfile) Name() string {
	if ip.ProfileName == "" {
		ip.Endpoint().DefaultProfile().Name()
	}

	return ip.ProfileName
}

func (ip ImageProfile) SetName(name string) profile.Profile {
	ip.ProfileName = name
	return ip
}

func (ip ImageProfile) ProfileRepository() profile.Repository {
	return profile.RuntimeRepository
}

func (ip ImageProfile) Endpoint() profile.Endpoint {
	return iEndpoint
}

func (ip *ImageProfile) Load(profileName string) (err error) {
	buf, err := profile.Repository.Read(ip.ProfileRepository(), profileName, ip.Endpoint().Name())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, ip)
	//if err != nil { return }

	return
}

// Use this function to generate the initial config for a profile
func getDefaultProfile() (profile ImageProfile, err error) {
	defaultProfileName, err := config.RuntimeConfig.GetDefaultProfile(ImageProfile{}.Endpoint().Name())
	if err != nil {
		return
	}

	err = profile.Load(defaultProfileName)
	if err != nil {
		return
	}

	return
}
