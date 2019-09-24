package config

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
	AppName string `json:"app_name"`
	Addr    string `json:"addr"`
	Domain 	string `json:"domain"`
	Redis   struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	Db struct {
		Host     string `json:"host"`
		Port     int 	`json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"db"`
	Auth struct {
		Github	struct {
			ClientId 		string 	`json:"client_id"`
			ClientSecret 	string 	`json:"client_secret"`
		}	`json:"github"`
	}	`json:"auth"`
	Salt 	string 	`json:"salt"`
}

func FindConfigFile() (string, error) {
	pwd, _ := os.Getwd()
	scanPath := []string{
		pwd,
		filepath.Dir(pwd),
	}
	for _, v := range scanPath {
		file := path.Join(v, "app.json")
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return file, nil
		}
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
		panic("config file read error " + err.Error())
	}

	err = json.Unmarshal(configData, &appConfig)
	if err != nil {
		panic("config file decode error " + err.Error())
	}

	return appConfig
}
