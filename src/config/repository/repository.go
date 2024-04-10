package repository

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/config/profile"
	"github.com/ephex2/go-gpt-cli/log"
)


// fileRepository implements both the ProfileRepository and ConfigRepository interfaces
type fileRepository struct {
	basePath        string
	filePath        string
	profileFileName string
}

func (cr *fileRepository) FilePath() (string, error) {
	return filepath.Abs(cr.filePath)
}

func (cr *fileRepository) Init() (err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}

	// Set default values
	cr.basePath = homeDir + "/.local/go-gpt-cli/"
	cr.filePath = cr.basePath + "go-gpt-cli.json"
	cr.profileFileName = "/config.json"

	// Make config folder
	_, err = os.Stat(cr.basePath)
	if err != nil {
		err = os.MkdirAll(cr.basePath, 0750)
		if err != nil && !os.IsExist(err) {
			return
		}
	}

	// Make empty config.json if it doesn't exist
	_, err = os.Stat(cr.filePath)
	if err != nil {
		_, err = os.Create(cr.filePath)
		if err != nil {
			return
		}
	}

    err = config.RuntimeConfig.Init(*cr)
    if err != nil {
        return
    }

	return
}

// ConfigRepository implementation
func (cr fileRepository) Get() (baseConfig config.Config, err error) {
	path, err := cr.FilePath()
	if err != nil {
		return
	}

	buf, err := os.ReadFile(path)
	if err != nil {
		return
	}

	settingsMap := make(map[string]string)
    
    if len(buf) != 0 {
        err = json.Unmarshal(buf, &settingsMap)
	    if err != nil {
		    return
	    }
    }

	baseConfig = config.Config{
		Repository: cr,
		Settings:   settingsMap,
	}

	return baseConfig, nil
}

func (cr fileRepository) Set(c config.Config) (err error) {
	path, err := cr.FilePath()
	if err != nil {
		return
	}

	settingsJson, err := json.Marshal(c.Settings)
	if err != nil {
		return
	}

	log.Debug("Writing config file to path: %s\n", path)
	err = os.WriteFile(path, settingsJson, 0750)
	if err != nil {
		return
	}

	return
}

// Implementing ProfileRepository interface directly on fileRepository
func (cr *fileRepository) profileFilePath(endpointName string, profileName string) string {
	return cr.profileFolderPath(endpointName, profileName) + cr.profileFileName
}

func (cr *fileRepository) profileFolderPath(endpointName string, profileName string) string {
	return cr.basePath + endpointName + "/" + profileName
}

func (cr fileRepository) Create(endpoint profile.Endpoint, profileName string) (err error) {
	p := endpoint.DefaultProfile()
	p = p.SetName(profileName)

	profileFolder := cr.profileFolderPath(endpoint.Name(), p.Name())
	err = os.MkdirAll(profileFolder, 0750)
	if err != nil && !os.IsExist(err) {
		panic(err.Error())
	} else {
		err = nil
	}

	err = cr.Update(p)
	if err != nil {
		return
	}

	err = config.SetDefaultProfile(endpoint.Name(), profileName, false)
    if err != nil {
        return
    }

	return
}

func (cr fileRepository) Read(name string, endpointName string) (pBytes []byte, err error) {
	log.Debug("Looking for profile in path: %s\n", cr.profileFilePath(endpointName, name))
	pBytes, err = os.ReadFile(cr.profileFilePath(endpointName, name))
	if err != nil {
		return
	}

	return
}

func (cr fileRepository) Update(p profile.Profile) (err error) {
	buf, err := json.Marshal(p)
	if err != nil {
		return
	}

	profilePath := cr.profileFilePath(p.Endpoint().Name(), p.Name())
	log.Debug("Writing profile at path: " + profilePath + "\n")
	err = os.WriteFile(cr.profileFilePath(p.Endpoint().Name(), p.Name()), buf, 0750)
	return
}

func (cr fileRepository) Delete(endpointName string, profileName string) (err error) {
	profilePath := cr.profileFolderPath(endpointName, profileName)
	log.Debug("Deleting profile at path: " + profilePath + "\n")
	err = os.RemoveAll(profilePath)
	if err != nil {
		return
	}

	if config.GetDefaultProfile(endpointName) == profileName {
		config.SetDefaultProfile(endpointName, "", true)
	}

	return
}

func (cr fileRepository) GetAll(endpointName string) (names []string, err error) {
	dirs, err := os.ReadDir(cr.basePath + "/" + endpointName)
	if err != nil || len(dirs) == 0 {
		return
	}

	for _, dir := range dirs {
		names = append(names, dir.Name())
	}

	return
}

// Why Init() ? The initialization of main seemed to only call init in packages imported by main.
// This function is meant to be called by the init function in the root cobra package
func Init() {
	dummyConfig := fileRepository{}
	err := dummyConfig.Init()
	if err != nil {
		panic(err.Error())
	}

	cfg, err := dummyConfig.Get()
	if err != nil {
		log.Warning("Error while getting global settings: " + err.Error() + "\n")
		log.Warning("\t~~~ During executions before any settings have been configured, this is normal ~~~~\t\n")
	}

	Profile := fileRepository{}
	err = Profile.Init()
	if err != nil {
		panic(err.Error())
	}

	config.RuntimeConfig = cfg
	config.RuntimeConfig.Repository = dummyConfig
	profile.RuntimeRepository = Profile
}
