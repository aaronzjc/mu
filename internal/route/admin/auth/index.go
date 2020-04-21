package auth

import (
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/route/middleware"
	"mu/internal/util/req"
)

func Info(c *gin.Context) {
	username, err := c.Cookie(middleware.CooAdmin)
	if err != nil {
		req.JSON(c, req.CodeError, "sorry, fetch admin failed", nil)
		return
	}

	if username == "" {
		req.JSON(c, req.CodeError, "empty user", nil)
	}

	admin := model.Admin{
		Username: username,
	}
	err = admin.FetchInfo()
	if err != nil {
		req.JSON(c, req.CodeError, "sorry, fetch user failed", nil)
		return
	}

	req.JSON(c, req.CodeSuccess, "userinfo", admin)
}
