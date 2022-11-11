package store

import (
	"context"
	"testing"

	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/test"

	"github.com/stretchr/testify/require"
)

func TestUserGetAll(t *testing.T) {
	require := require.New(t)

	test.SetupTestDb(t, model.DB_MU)

	repo := NewUserRepo()

	users, err := repo.GetAll(context.Background())
	require.Nil(err)
	require.NotEmpty(users)
}
