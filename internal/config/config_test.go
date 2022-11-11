package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupConfig() string {
	f, _ := os.CreateTemp("", "config")
	f.WriteString("name: api\nenv: dev\nhost: 127.0.0.1\nport: 8780")
	return f.Name()
}

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)

	conf, err := LoadConfig("invalid path")
	assert.NotNil(err)
	assert.Nil(conf)

	f := setupConfig()
	defer os.Remove(f)
	conf, err = LoadConfig(f)
	assert.Nil(err)
	assert.Equal(conf.Http.Port, 8780)
}
