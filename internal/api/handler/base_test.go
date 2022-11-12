package handler

import (
	"testing"

	"github.com/aaronzjc/mu/internal/constant"
	"github.com/aaronzjc/mu/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResp(t *testing.T) {
	assert := assert.New(t)
	h := func(c *gin.Context) {
		Resp(c, constant.CodeSuccess, "ok", nil)
	}

	resp := test.NewRequest(t).Handler(h).Get("/").Exec()
	assert.Equal(200, resp.Code())
	errno, errmsg, _, _ := resp.TryDecode()
	assert.Equal(errno, constant.CodeSuccess)
	assert.Equal(errmsg, "ok")
}
