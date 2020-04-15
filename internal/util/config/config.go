package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var appConfig Config

type Config struct {
	AppName string `json:"app_name"`
	Server  struct {
		Https bool   `json:"https"`
		Host  string `json:"host"`
		Addr  string `json:"addr"`
		Web   string `json:"web"`
	} `json:"server"`
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	Db struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"db"`
	Auth struct {
		Github struct {
			ClientId     string   `json:"client_id"`
			ClientSecret string   `json:"client_secret"`
			Admins       []string `json:"admins"`
		} `json:"github"`
	} `json:"auth"`
	Salt string `json:"salt"`
}

func (c *Config) ServerUrl() string {
	proto := "http"
	if c.Server.Https {
		proto = "https"
	}

	return fmt.Sprintf("%s://%s", proto, c.Server.Host)
}

func (c *Config) WebUrl() string {
	return c.Server.Web
}

func FindConfigFile() (string, error) {
	pwd, _ := os.Getwd()
	scanPath := []string{
		pwd,
		filepath.Dir(pwd),
	}
	for _, v := range scanPath {
		file := path.Join(v, "/conf/app.json")
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return file, nil
		}
	}

	return "", errors.New("couldn't find config file")
}

func NewConfig() Config {
	if appConfig.AppName != "" {
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

	if appConfig.AppName == "" {
		panic("invalid config")
	}

	// 如果前端Web为空，则表明使用Gin的Web服务。这里，更新Web地址
	if appConfig.Server.Web == "" {
		appConfig.Server.Web = appConfig.ServerUrl()
	}

	return appConfig
}
