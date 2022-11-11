package config

import (
	"errors"

	"github.com/aaronzjc/mu/internal/constant"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type HttpConfig struct {
	Tls  bool   `yaml:"tls"`
	Url  string `yaml:"url"`
	Port int    `yaml:"port"`
}

type CommanderConfig struct {
	Port int `yaml:"port"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type ServiceConfig struct {
	Url string `yaml:"url"`
}

type OAuthConfig struct {
	ClientId     string   `yaml:"clientId"`
	ClientSecret string   `yaml:"clientSecret"`
	Admins       []string `yaml:"admins"`
}

type Config struct {
	Name      string                   `yaml:"name"`
	Env       string                   `yaml:"env"`
	Salt      string                   `yaml:"salt"`
	Log       LogConfig                `yaml:"log"`
	Http      HttpConfig               `yaml:"http"`
	Commander CommanderConfig          `yaml:"commander"`
	Redis     RedisConfig              `yaml:"redis"`
	Database  map[string]DbConfig      `yaml:"database"`
	Service   map[string]ServiceConfig `yaml:"service"`
	OAuth     map[string]OAuthConfig   `yaml:"oauth"`
}

func (c *Config) GetServiceUrl(name constant.SvcName) string {
	svc, ok := c.Service[string(name)]
	if !ok {
		return ""
	}
	return svc.Url
}

func (c *Config) ServerUrl() string {
	if c.Http.Tls {
		return "https://" + c.Http.Url
	}
	return "http://" + c.Http.Url
}

func (c *Config) AdminUrl() string {
	return c.ServerUrl() + "/admin#/"
}

func (c *Config) IndexUrl() string {
	return c.ServerUrl() + "/#/"
}

var (
	vip    *viper.Viper
	config *Config
)

func init() {
	vip = viper.New()
	config = &Config{}
}

func LoadConfig(path string) (*Config, error) {
	//加载配置
	vip.SetConfigFile(path)
	vip.SetConfigType("yml")
	if err := vip.ReadInConfig(); err != nil {
		return nil, errors.New("read config file error")
	}
	vip.Unmarshal(&config)

	// 监听配置变更
	vip.OnConfigChange(func(in fsnotify.Event) {
		// do something
	})
	vip.WatchConfig()

	return config, nil
}

func Get() *Config {
	return config
}
