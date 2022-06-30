package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"

	"github.com/docker/docker-credential-helpers/credentials"
)

var ErrNoCredentials = errors.New("missing credential informations")

type Secret struct {
	Auth     string `json:"auth,omitempty"`
	Username string `json:"username,omitempty"`
	Secret   string `json:"secret,omitempty"`
}

func (s *Secret) ToUserPass() error {
	if s.Auth == "" && s.Username == "" && s.Secret == "" {
		return ErrNoCredentials
	}

	decoded, err := base64.StdEncoding.DecodeString(s.Auth)
	if err != nil {
		return err
	}

	split := strings.Split(string(decoded), ":")
	if len(split) != 2 {
		return fmt.Errorf("cannot decode: %s", decoded)
	}

	s.Username = split[0]
	s.Secret = split[1]
	s.Auth = ""

	return nil
}

func (s *Secret) ToAuth() error {
	if s.Username == "" || s.Secret == "" {
		return ErrNoCredentials
	}

	s.Auth = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Username, s.Secret)))
	s.Username = ""
	s.Secret = ""

	return nil
}

func FromCredentials(creds *credentials.Credentials) *Secret {
	return &Secret{
		Username: creds.Username,
		Secret:   creds.Secret,
	}
}

type FileBased struct {
	path string
}

func NewFileBasedStore() (credentials.Helper, error) {
	path, err := GetConfPath()
	if err != nil {
		return nil, err
	}

	return newFileBasedStore(path)
}

func newFileBasedStore(path string) (credentials.Helper, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", path)
	}

	return &FileBased{
		path: path,
	}, nil
}

func (f *FileBased) Add(creds *credentials.Credentials) error {
	data, err := readData(f.path)
	if err != nil {
		return err
	}

	newSecret := FromCredentials(creds)
	if err := newSecret.ToAuth(); err != nil {
		return err
	}

	var auths map[string]interface{}
	if authsVal, ok := data["auths"].(map[string]interface{}); !ok {
		auths = make(map[string]interface{})
	} else {
		auths = authsVal
	}

	auths[creds.ServerURL] = map[string]string{
		"auth": newSecret.Auth,
	}

	data["auths"] = auths

	return writeData(f.path, data)
}

func (f *FileBased) Delete(serverURL string) error {
	data, err := readData(f.path)
	if err != nil {
		return err
	}

	auths, ok := data["auths"].(map[string]interface{})
	if !ok {
		return ErrNoCredentials
	}

	delete(auths, serverURL)

	data["auths"] = auths

	return writeData(f.path, data)
}

func (f *FileBased) Get(serverURL string) (string, string, error) {
	data, err := readData(f.path)
	if err != nil {
		return "", "", err
	}

	auths, ok := data["auths"].(map[string]interface{})
	if !ok {
		return "", "", ErrNoCredentials
	}

	creds, ok := auths[serverURL].(map[string]interface{})
	if !ok {
		return "", "", ErrNoCredentials
	}

	secret := &Secret{
		Auth: creds["auth"].(string),
	}
	if err := secret.ToUserPass(); err != nil {
		return "", "", err
	}

	return secret.Username, secret.Secret, nil
}

func (f *FileBased) List() (map[string]string, error) {
	data, err := readData(f.path)
	if err != nil {
		return nil, err
	}

	auths, ok := data["auths"].(map[string]interface{})
	if !ok {
		return nil, ErrNoCredentials
	}

	result := make(map[string]string)
	for serverURL, creds := range auths {
		c, ok := creds.(map[string]interface{})
		if !ok {
			continue
		}

		secret := &Secret{}
		if auth, ok := c["auth"].(string); ok {
			secret.Auth = auth
			if err := secret.ToUserPass(); err != nil {
				continue
			}
		} else {
			user, ok := c["username"].(string)
			if !ok {
				continue
			}
			secret.Username = user
		}

		result[serverURL] = secret.Username
	}

	return result, nil
}

func readData(path string) (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func writeData(path string, data map[string]interface{}) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, out, 0600)
}

func GetConfPath() (string, error) {
	var path string

	env := os.Getenv("DOCKER_CONFIG")
	if env != "" {
		path = env
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, ".docker", "config.json")
	}

	return path, nil
}
