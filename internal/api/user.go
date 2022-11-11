package api

import (
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/gin-gonic/gin"
)

type User struct {
	svc service.UserService
}

func (u *User) List(ctx *gin.Context) {
	users, err := u.svc.GetUserList(ctx)
	if err != nil {
		Resp(ctx, constant.CodeError, err.Error(), nil)
		return
	}
	Resp(ctx, constant.CodeSuccess, "success", users)
}

func NewUser() *User {
	repo := store.NewUserRepo()
	svc := service.NewUserService(repo)
	return &User{
		svc: svc,
	}
}
