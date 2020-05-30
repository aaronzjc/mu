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
		Https  bool   `json:"https"`
		Host   string `json:"host"`
		Port   string `json:"port"`
		Static bool   `json:"static"`
	} `json:"server"`
	Frontend struct {
		Index string `json:"index"`
		Admin string `json:"admin"`
	} `json:"frontend"`
	Commander struct {
		Addr string `json:"addr"`
	} `json:"commander"`
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

// 如果后端托管静态资源，则构造路由。
// 因为和api公用域名，为了不混淆，不支持history模式
func (c *Config) IndexUrl() string {
	if c.Server.Static {
		return c.ServerUrl() + "/#/"
	}
	if c.Frontend.Index == "" {
		panic("index route empty")
	}
	return c.Frontend.Index
}

// 如果后端托管静态资源，则构造路由。
// 因为和api公用域名，为了不混淆，不支持history模式
func (c *Config) AdminUrl() string {
	if c.Server.Static {
		return c.ServerUrl() + "/admin#/"
	}
	if c.Frontend.Admin == "" {
		panic("admin route empty")
	}
	return c.Frontend.Admin
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

	return appConfig
}
