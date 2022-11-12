package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	svc    service.OAuthService
	usrSvc service.UserService
}

func (o *Auth) LoginInfo(ctx *gin.Context) {
	userId, ok := ctx.Get(constant.LoginKey)
	if !ok {
		Resp(ctx, constant.CodeError, "用户不存在", nil)
		return
	}
	user, _ := o.usrSvc.GetUser(ctx, &dto.Query{
		Query: "`id` = ?",
		Args:  []interface{}{userId.(int)},
	})
	Resp(ctx, constant.CodeSuccess, "success", user)
}

func (o *Auth) Platforms(ctx *gin.Context) {
	from := ctx.Query("from")
	if from != "admin" {
		from = "index"
	}
	platforms := o.svc.Platforms(from)
	Resp(ctx, constant.CodeSuccess, "success", platforms)
}

func (o *Auth) Auth(ctx *gin.Context) {
	from := ctx.Query("from")
	by := ctx.Query("by")

	// 设置登录后跳转
	SetCookies(ctx, map[string]string{
		"from": from,
		"by":   by,
	}, "")

	redirect := o.svc.Redirect(by)
	if redirect == "" {
		ctx.Abort()
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, redirect)
	ctx.Abort()
}

func (o *Auth) Callback(ctx *gin.Context) {
	var from, by string
	var err error
	code := ctx.Query("code")
	from, err = ctx.Cookie("from")
	if err != nil || from != "admin" {
		// 默认跳首页
		from = "index"
	}
	by, err = ctx.Cookie("by")
	if err != nil {
		ctx.String(http.StatusForbidden, "get cookie err")
		return
	}

	if code == "" {
		ctx.String(http.StatusForbidden, "code缺失")
		return
	}

	oauthUser, err := o.svc.Auth(by, code)
	if err != nil {
		ctx.String(http.StatusForbidden, err.Error())
		return
	}
	token, err := o.usrSvc.Auth(ctx, by, oauthUser)
	if err != nil {
		ctx.String(http.StatusForbidden, err.Error())
		return
	}
	tokenEncode := url.QueryEscape(token)
	cnf := config.Get()

	var redirect = ""
	if from == "admin" {
		redirect = fmt.Sprintf("%s?token=%s", cnf.AdminUrl(), tokenEncode)
	} else {
		redirect = fmt.Sprintf("%s?token=%s", cnf.IndexUrl(), tokenEncode)
	}

	ctx.Redirect(http.StatusTemporaryRedirect, redirect)
	ctx.Abort()
}

func NewAuth() *Auth {
	return &Auth{
		svc:    service.NewOAuthService(),
		usrSvc: service.NewUserService(store.NewUserRepo()),
	}
}
