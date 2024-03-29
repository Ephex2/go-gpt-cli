package config

import (
	"errors"
)

var SetKeyDoc = "go-gpt-cli setKey abc123"

func GetApiKey() (string, error) {
	if value, ok := RuntimeConfig.Settings["ApiKey"]; ok {
		return value, nil
	} else {
		return "", errors.New("ApiKey not defined in settings. Set it using SetApiKey, ex: " + SetKeyDoc)
	}
}

func SetApiKey(key string) (err error) {
	if RuntimeConfig.Settings == nil {
		m := make(map[string]string)
		RuntimeConfig.Settings = m
	}

	RuntimeConfig.Settings["ApiKey"] = key
	return RuntimeConfig.Repository.Set(RuntimeConfig)

	// no need to refresh since memory value = disk calue in happy path
	// error path will need to be determined by callerFile either terminate or refresh and continue.
}
