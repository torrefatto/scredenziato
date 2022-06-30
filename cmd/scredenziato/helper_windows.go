// go:build windows
package main

import (
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/docker/docker-credential-helpers/wincred"

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
	case "desktop", "wincred":
		return wincred.Wincred{}
	}
	return nil
}

func getChainHelper() (credentials.Helper, error) {
	helper := make(map[string]credentials.Helper)
	helper["wincred"] = wincred.Wincred{}
	fileHelper, err := helpers.NewFileBasedStore()
	if err != nil {
		return nil, err
	}
	helper["file"] = fileHelper

	return helpers.ChainHelper(helper), nil
}
