package config

import (
	"errors"
	"net/url"

	"github.com/ephex2/go-gpt-cli/config/profile"
)

//  This will be the base url for the config store if it is not set yet.
//  After a first run, the command SetBaseUrl of the config repository can be used to set this value
const defaultBaseUrl string = "https://api.openai.com"
const baseUrlKeyName string = "BaseUrl"

type ConfigRepository interface {
	Get() (Config, error)
	Set(Config) error
}

type Config struct {
	Settings   map[string]string
	Repository ConfigRepository
}

// The config to refer to at runtime. It contains settings that can be referenced by all endpoints as well as its own repository.
// This repository can be used to make modifications to the config.
var RuntimeConfig Config

func BaseUrl() string {
    if url, ok := RuntimeConfig.Settings[baseUrlKeyName] ; ok {
        return url
    } 

    return defaultBaseUrl
}

func SetBaseUrl(baseurl string) (err error) {
    u, err := url.Parse(baseurl)
	if err != nil {
		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New("Invalid URL, scheme is not http or https. It is: " + u.Scheme)
        return 
	}

	if u.Host == "" {
		err = errors.New("Invalid URL. URL provided had no host.")
        return
	}

    // Url ok, set config
    RuntimeConfig.Settings[baseUrlKeyName] = baseurl
    err = RuntimeConfig.Repository.Set(RuntimeConfig)
    if err != nil {
        return
    }

    return
}


// Used to initialize the config's settings stored in its repository
func (c *Config) Init(cr ConfigRepository) (err error) {
	*c, err = cr.Get()
	if err != nil {
		return
	}

    // If default URL isn't set, create 
    _, ok := c.Settings[baseUrlKeyName]
    if !ok {
        c.Settings[baseUrlKeyName] = defaultBaseUrl
        cr.Set(*c)
    }

	return
}

// Gets the default profile name associated with a given endpoint.
// If none exist, creates the default profile associated with that endpoint.
// This avoids errors when running the tool for the first time.
func (c *Config) GetDefaultProfile(endpointName string) (s string, err error) {
	s, ok := c.Settings[endpointName+"DefaultProfile"]

	// if default profile doesn't exist, create it
	if !ok || s == "" {
		var e profile.Endpoint
		e, err = profile.EndpointRegistry.Get(endpointName)
		if err != nil {
			return
		}

		err = profile.RuntimeRepository.Create(e, e.DefaultProfile().Name())
		if err != nil {
			return
		}

		s = e.DefaultProfile().Name()
		err = SetDefaultProfile(e.Name(), s, true)
	}

	return s, err
}

// Mixed on where this should be in the code, using KISS right now and leaving it in top-level Config.
// Sets default profile for endpoint if not already set. Force flag skips check.
func SetDefaultProfile(endpointName string, profileName string, force bool) (err error) {
	// Test if endpoint name exists
	_, err = profile.EndpointRegistry.Get(endpointName)
	if err != nil {
		return
	}

	retrievedName, ok := RuntimeConfig.Settings[endpointName+"DefaultProfile"]

	if !ok || force || retrievedName == "" {
		RuntimeConfig.Settings[endpointName+"DefaultProfile"] = profileName
		err = RuntimeConfig.Repository.Set(RuntimeConfig)
		if err != nil {
			return
		}
	}

	return
}

// Gets default profile, returns empty string if not found
func GetDefaultProfile(endpointName string) string {
	profileName := RuntimeConfig.Settings[endpointName+"DefaultProfile"]
	return profileName
}
