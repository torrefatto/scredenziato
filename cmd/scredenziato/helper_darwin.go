// go:build darwin
package main

import (
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/docker/docker-credential-helpers/osxkeychain"
	"github.com/docker/docker-credential-helpers/pass"

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
	case "desktop", "osxkeychain":
		return osxkeychain.Osxkeychain{}
	case "pass":
		return pass.Pass{}
	}
	return nil
}

func getChainHelper() (credentials.Helper, error) {
	helper := make(map[string]credentials.Helper)
	helper["osxkeychain"] = osxkeychain.Osxkeychain{}
	helper["pass"] = pass.Pass{}
	fileHelper, err := helpers.NewFileBasedStore()
	if err != nil {
		return nil, err
	}
	helper["file"] = fileHelper

	return helpers.ChainHelper(helper), nil
}
