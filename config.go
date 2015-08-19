package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"gopkg.in/yaml.v1"

	"github.com/phrase/phraseapp-go/phraseapp"
)

const configName = ".phraseapp.yml"
const defaultDir = "./"

type Credentials struct {
	phraseapp.Credentials
	Username string `cli:"opt --username -u desc='username used for authentication'"`
	Token    string `cli:"opt --access-token -t desc='access token used for authentication'"`
	TFA      bool   `cli:"opt --tfa desc='use Two-Factor Authentication'"`
	Host     string `cli:"opt --host desc='Host to send Request to'"`
	Debug    bool   `cli:"opt --verbose -v desc='Verbose output'"`
}

func ConfigCallArgs() (map[string]string, error) {
	content, err := ConfigContent()
	if err != nil {
		content = "{}"
	}

	return parseCallArgs(content)
}

func parseCallArgs(yml string) (map[string]string, error) {
	var callArgs *CallArgs

	err := yaml.Unmarshal([]byte(yml), &callArgs)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)

	if callArgs != nil {
		m["ProjectID"] = callArgs.Phraseapp.ProjectID
		m["AccessToken"] = callArgs.Phraseapp.AccessToken
	}

	return m, nil
}

func ClientFromCmdCredentials(c Credentials) (*phraseapp.Client, error) {
	defaultCredentials, e := ConfigDefaultCredentials()
	if e != nil {
		return nil, e
	}

	return phraseapp.NewClient(PhraseAppCredentials(c), defaultCredentials)
}

func PhraseAppCredentials(c Credentials) phraseapp.Credentials {
	return phraseapp.Credentials{
		Username: c.Username,
		Token:    c.Token,
		TFA:      c.TFA,
		Host:     c.Host,
		Debug:    c.Debug,
	}
}

func ConfigDefaultCredentials() (*phraseapp.Credentials, error) {
	content, err := ConfigContent()
	if err != nil {
		content = "{}"
	}

	return parseCredentials(content)
}

func ConfigDefaultParams() (phraseapp.DefaultParams, error) {
	content, err := ConfigContent()
	if err != nil {
		content = "{}"
	}

	return parseDefaults(content)
}

func ConfigContent() (string, error) {
	path, err := configPath()
	if err != nil {
		return "", err
	}

	bytes, err := readFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func configPath() (string, error) {
	if envConfig := os.Getenv("PHRASEAPP_CONFIG"); envConfig != "" {
		possiblePath := path.Join(envConfig)
		if _, err := os.Stat(possiblePath); err == nil {
			return possiblePath, nil
		}
	}

	callerPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	possiblePath := path.Join(callerPath, configName)
	if _, err := os.Stat(possiblePath); err == nil {
		return possiblePath, nil
	}

	return defaultConfigDir()
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func defaultConfigDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", nil
	}
	return path.Join(usr.HomeDir, defaultDir, configName), nil
}

// Parsing
type credentialConf struct {
	Phraseapp struct {
		AccessToken string `yaml:"access_token"`
		Host        string `yaml:"host"`
		Debug       bool   `yaml:"verbose"`
		Username    string
		TFA         bool
	}
}

func parseCredentials(yml string) (*phraseapp.Credentials, error) {
	var conf *credentialConf

	if err := yaml.Unmarshal([]byte(yml), &conf); err != nil {
		fmt.Fprintln(os.Stderr, "Could not parse .phraseapp.yml")
		return nil, err
	}

	phrase := conf.Phraseapp

	credentials := &phraseapp.Credentials{Token: phrase.AccessToken, Username: phrase.Username, TFA: phrase.TFA, Host: phrase.Host, Debug: phrase.Debug}

	return credentials, nil
}

type defaultsConf struct {
	Phraseapp struct {
		Defaults phraseapp.DefaultParams
	}
}

func parseDefaults(yml string) (phraseapp.DefaultParams, error) {
	var conf *defaultsConf

	err := yaml.Unmarshal([]byte(yml), &conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not parse .phraseapp.yml")
		return nil, err
	}

	return conf.Phraseapp.Defaults, nil
}

type CallArgs struct {
	Phraseapp struct {
		AccessToken string `yaml:"access_token"`
		ProjectID   string `yaml:"project_id"`
		Page        int
		PerPage     int
	}
}
