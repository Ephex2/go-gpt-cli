package batches

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/config/profile"
)

type BatchProfile struct {
	ProfileName    string
	CreateBatchBody CreateBatchBody
    Url            string
}

func (p BatchProfile) Name() string {
	if p.ProfileName == "" {
		p.Endpoint().DefaultProfile().Name()
	}

	return p.ProfileName
}

func (p BatchProfile) OverrideUrl() string {
    return p.Url
}

func (p BatchProfile) SetName(name string) profile.Profile {
	p.ProfileName = name
	return p
}

func (p BatchProfile) ProfileRepository() profile.Repository {
	return profile.RuntimeRepository
}

func (p BatchProfile) Endpoint() profile.Endpoint {
	return bEndpoint
}

func (p BatchProfile) Validate() error {
    var err error
    err = AllowedCompletionWindow(p.CreateBatchBody.CompletionWindow)
    err = AllowedBatchApiEndpoint(p.CreateBatchBody.Endpoint)
    
    return err
}

func (p *BatchProfile) Load(profileName string) (err error) {
	buf, err := profile.Repository.Read(p.ProfileRepository(), profileName, p.Endpoint().Name())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, p)
	//if err != nil { return }

	return
}

// Use this function to generate the initial config for a profile
func getDefaultProfile() (profile BatchProfile, err error) {
	defaultProfileName, err := config.RuntimeConfig.GetDefaultProfile(BatchProfile{}.Endpoint().Name())
	if err != nil {
		return
	}

	err = profile.Load(defaultProfileName)
	if err != nil {
		return
	}

	return
}
