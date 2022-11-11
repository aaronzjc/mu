package db

import (
	"testing"

	"github.com/aaronzjc/mu/internal/config"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	conf = config.Config{
		Database: map[string]config.DbConfig{
			"demo": {
				Host:     "127.0.0.1",
				Port:     3306,
				Username: "root",
				Password: "123456",
				Charset:  "utf8",
			},
		},
	}
)

func TestDb(t *testing.T) {
	require := require.New(t)

	err := Setup(&conf, &gorm.Config{})
	require.Nil(err)

	demo, ok := Get("demo")
	require.True(ok)
	require.NotEmpty(demo)

	db, _ := demo.DB()
	require.Nil(db.Ping())
}
