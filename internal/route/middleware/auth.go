package middleware

import (
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/aaronzjc/mu/internal/api"
	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/gin-gonic/gin"
)

func ApiAuth(admin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			api.Resp(ctx, constant.CodeAuthFailed, "认证失败", nil)
			ctx.Abort()
			return
		}
		token, _ = url.QueryUnescape(token)
		authBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			api.Resp(ctx, constant.CodeAuthFailed, "解析失败", nil)
			ctx.Abort()
			return
		}
		segs := strings.Split(string(authBytes), ";")
		if len(segs) != 2 || segs[0] == "" || segs[1] == "" {
			api.Resp(ctx, constant.CodeAuthFailed, "token格式有误", nil)
			ctx.Abort()
			return
		}

		username, token := segs[0], segs[1]

		userRepo := store.NewUserRepo()
		userService := service.NewUserService(userRepo)

		if ok := userService.VerifyToken(ctx, username, token); !ok {
			api.Resp(ctx, constant.CodeAuthFailed, "禁止访问", nil)
			ctx.Abort()
			return
		}

		u, _ := userService.GetUser(ctx, &dto.Query{
			Query: "`username` = ?",
			Args:  []interface{}{username},
		})

		if admin {
			admins := []string{"everyone"}
			oauthConfig, ok := config.Get().OAuth[u.AuthType]
			if ok {
				admins = oauthConfig.Admins
			}
			for _, v := range admins {
				if u.Username != v && v != "everyone" {
					api.Resp(ctx, constant.CodeForbidden, "不好意思，您没有权限。请联系管理员", nil)
					ctx.Abort()
					return
				}
			}
		}
		ctx.Set(constant.LoginKey, u.ID)
		ctx.Next()
	}
}
