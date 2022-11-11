package test

import (
	"testing"

	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/infra/db"

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

func SetupTestDb(t *testing.T, dbName string) {
	_, ok := db.Get(dbName)
	if !ok {
		require.Nil(t, db.Setup(&conf, &gorm.Config{}))
	}
}
