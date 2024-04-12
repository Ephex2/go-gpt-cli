package audio

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/config/profile"
)

type AudioProfile struct {
	ProfileName                    string
	CreateSpeechBody               CreateSpeechBody
	CreateTranscriptionBody        map[string]string
	CreateVerboseTranscriptionBody map[string]string
	CreateTranslationBody          map[string]string
	CreateVerboseTranslationBody   map[string]string
	SaveDirectory                  string // Blank by default, results in temporary files being created if blank
    Url                            string
}

func (a AudioProfile) Name() string {
	if a.ProfileName == "" {
		a.Endpoint().DefaultProfile().Name()
	}

	return a.ProfileName
}

func (a AudioProfile) OverrideUrl() string {
    return a.Url
}

func (a AudioProfile) SetName(name string) profile.Profile {
	a.ProfileName = name
	return a
}

func (a AudioProfile) ProfileRepository() profile.Repository {
	return profile.RuntimeRepository
}

func (a AudioProfile) Endpoint() profile.Endpoint {
	return aEndpoint
}

func (a *AudioProfile) Load(profileName string) (err error) {
	buf, err := profile.Repository.Read(a.ProfileRepository(), profileName, a.Endpoint().Name())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, a)
	//if err != nil { return }

	return
}

// Saves a file with name "name" to the profile's SaveFolder. Returns the filePath only if the file was successfully written.
func (a AudioProfile) Save(buf []byte, ext string) (filePath string, err error) {
	var f *os.File

	if a.SaveDirectory == "" {
		f, err = os.CreateTemp("", "audioFile-*"+ext)
	} else {
		f, err = os.CreateTemp(a.SaveDirectory, "audioFile-*"+ext)
	}

	if err != nil {
		return
	}
	defer f.Close()

	rc := bytes.NewReader(buf)
	_, err = io.Copy(f, rc)
	if err != nil {
		return
	}

	filePath = f.Name()
	return
}

// Use this function to generate the initial config for a profile
func getDefaultProfile() (profile AudioProfile, err error) {
	defaultProfileName, err := config.RuntimeConfig.GetDefaultProfile(AudioProfile{}.Endpoint().Name())
	if err != nil {
		return
	}

	err = profile.Load(defaultProfileName)
	if err != nil {
		return
	}

	return
}
