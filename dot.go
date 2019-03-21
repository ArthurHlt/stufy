package stufy

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
)

const filenameDot string = ".stufy"

type DotConfig struct {
	Aliases Aliases `json:"aliases"`
}

type Aliases []Alias

func (as Aliases) FindTarget(name string) string {
	for _, a := range as {
		if a.Name == name {
			return a.Target
		}
	}
	return ""
}

type Alias struct {
	Name   string `json:"name"`
	Target string `json:"target"`
}

func loadDotConfig() (*DotConfig, error) {
	confPath, err := dotConfPath()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(confPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err != nil && os.IsNotExist(err) {
		emptyFile, err := os.Create(confPath)
		if err != nil {
			return nil, err
		}
		emptyFile.WriteString("{}")
		emptyFile.Close()
		return &DotConfig{}, nil
	}
	dotConfig := &DotConfig{}
	err = json.Unmarshal(b, dotConfig)
	if err != nil {
		return nil, err
	}
	return dotConfig, nil
}

func dotConfPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	confPath := filepath.Join(home, filenameDot)
	return confPath, nil
}

func saveDotConfig(dotConfig *DotConfig) error {
	confPath, err := dotConfPath()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(dotConfig, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(confPath, b, 0644)
}
