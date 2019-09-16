package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var appConfig Config

type Config struct {
	AppName		string 		`json:"app_name"`
	Redis 		struct{
		Host 		string		`json:"host"`
		Port		int			`json:"port"`
	}
}

func FindConfigFile() (string, error) {
	pwd, _ := os.Getwd()
	currPath := path.Join(pwd, "app.json")
	if _, err := os.Stat(currPath); !os.IsNotExist(err) {
		return currPath, nil
	}
	currPath = path.Join(filepath.Dir(pwd), "app.json")
	if _, err := os.Stat(currPath); !os.IsNotExist(err) {
		return currPath, nil
	}

	return "", errors.New("couldn't find config file")
}

func NewConfig() Config {
	if appConfig != (Config{}) {
		return appConfig
	}

	var err error
	configPath, err := FindConfigFile()
	if err != nil {
		panic(err)
	}
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic("config file read error")
	}

	err = json.Unmarshal(configData, &appConfig)
	if err != nil {
		panic("config file decode error")
	}

	return appConfig
}