package helpers

import (
	"github.com/docker/docker-credential-helpers/credentials"
)

type ChainHelper map[string]credentials.Helper

var Chain credentials.Helper = make(ChainHelper)

func (c ChainHelper) Add(creds *credentials.Credentials) error {
	var errs ChainErr

	for name, helper := range c {
		if err := helper.Add(creds); err == nil {
			return nil
		} else {
			errs.Add(name, err)
		}
	}

	return errs
}

func (c ChainHelper) Delete(serverURL string) error {
	errs := NewChainErr()

	for name, helper := range c {
		if err := helper.Delete(serverURL); err == nil {
			return nil
		} else {
			errs.Add(name, err)
		}
	}

	return errs
}

func (c ChainHelper) Get(serverURL string) (string, string, error) {
	errs := NewChainErr()

	for name, helper := range c {
		username, secret, err := helper.Get(serverURL)
		if err == nil {
			return username, secret, nil
		}
		errs.Add(name, err)
	}

	return "", "", errs
}

func (c ChainHelper) List() (map[string]string, error) {
	errs := NewChainErr()

	for name, helper := range c {
		list, err := helper.List()
		if err == nil {
			return list, nil
		}

		errs.Add(name, err)
	}

	return nil, errs
}
