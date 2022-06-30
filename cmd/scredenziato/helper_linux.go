// go:build linux
package main

import (
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/docker/docker-credential-helpers/pass"
	"github.com/docker/docker-credential-helpers/secretservice"

	"github.com/torrefatto/scredenziato/helpers"
)

func getHelper() (credentials.Helper, error) {
	if store := getHelperFromFile(); store != nil {
		return store, nil
	}
	return getChainHelper()
}

func tryOsSpecificHelpers(credStore string) credentials.Helper {
	switch credStore {
	case "secretservice":
		return secretservice.Secretservice{}
	case "pass":
		return pass.Pass{}
	}
	return nil
}

func getChainHelper() (credentials.Helper, error) {
	helper := make(map[string]credentials.Helper)
	helper["secretservice"] = secretservice.Secretservice{}
	helper["pass"] = pass.Pass{}
	fileHelper, err := helpers.NewFileBasedStore()
	if err != nil {
		return nil, err
	}
	helper["file"] = fileHelper

	return helpers.ChainHelper(helper), nil
}
