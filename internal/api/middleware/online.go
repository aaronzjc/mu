package middleware

import (
	"fmt"
	"net/http"

	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/aaronzjc/mu/pkg/helper"
	"github.com/gin-gonic/gin"
)

func SetOnline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer ctx.Next()
		clientIp := helper.ClientIp(ctx.Request)
		if clientIp == "" {
			return
		}
		svcUrl := config.Get().GetServiceUrl(constant.SvcOnline)
		if svcUrl == "" {
			return
		}
		url := fmt.Sprintf("%s/online/%s/%s", svcUrl, "mu", clientIp)
		req, _ := http.NewRequest("POST", url, nil)
		resp, err := (new(http.Client)).Do(req)
		if err == nil {
			resp.Body.Close()
		}
	}
}
