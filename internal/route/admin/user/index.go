package user

import (
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/util/req"
)

func List(c *gin.Context) {
	list, err := (&model.User{}).FetchRows(model.Query{})
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}

	var result []model.UserJson
	for _, val := range list {
		u, _ := val.FormatJson()
		result = append(result, u)
	}
	req.JSON(c, req.CodeSuccess, "成功", result)
}
