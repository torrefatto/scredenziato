//go:build linux

package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/docker/docker-credential-helpers/credentials"

	"github.com/torrefatto/scredenziato/helpers"
)

func getHelperFromFile() credentials.Helper {
	path, err := helpers.GetConfPath()
	if err != nil {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil
	}

	if credStore, ok := config["credsStore"].(string); ok {
		return tryOsSpecificHelpers(credStore)
	} else {
		if _, ok := config["auths"].(map[string]interface{}); ok {
			fileHelper, err := helpers.NewFileBasedStore()
			if err != nil {
				return nil
			}
			return fileHelper
		}
	}

	return nil
}
