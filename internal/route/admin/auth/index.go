package auth

import (
	"crawler/internal/model"
	"crawler/internal/route/middleware"
	"crawler/internal/util/req"
	"github.com/gin-gonic/gin"
)

func Info(c *gin.Context) {
	username, err := c.Cookie(middleware.CooAdmin)
	if err != nil {
		req.JSON(c, req.CodeError, "sorry, fetch admin failed", nil)
		return
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
