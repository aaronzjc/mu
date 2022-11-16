package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/gin-gonic/gin"
)

type Stat struct{}

func (ctr *Stat) Online(c *gin.Context) {
	svcUrl := config.Get().GetServiceUrl(constant.SvcOnline)
	if svcUrl == "" {
		Resp(c, constant.CodeSuccess, "success", map[string]string{
			"count": "",
		})
		return
	}
	url := fmt.Sprintf("%s/online/%s", svcUrl, "mu")
	resp, err := http.Get(url)
	if err != nil {
		Resp(c, constant.CodeSuccess, "success", map[string]string{
			"count": "",
		})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Resp(c, constant.CodeSuccess, "success", map[string]string{
			"count": "",
		})
		return
	}
	Resp(c, constant.CodeSuccess, "success", map[string]string{
		"count": string(body),
	})
}

func NewStat() *Stat {
	return &Stat{}
}
