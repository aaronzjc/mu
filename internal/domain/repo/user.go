package repo

import (
	"context"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
)

// User 用户相关行为
type UserRepo interface {
	GetUsers(context.Context, *dto.Query) ([]model.User, error)
	GetUser(context.Context, *dto.Query) (model.User, error)
	CreateUser(context.Context, model.User) error
	Update(context.Context, model.User, map[string]interface{}) error
}
