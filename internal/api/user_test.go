package api

import (
	"testing"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/aaronzjc/mu/internal/mocks"
	"github.com/aaronzjc/mu/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mockUser(t *testing.T) *User {
	svc := mocks.NewUserService(t)
	svc.EXPECT().GetUserList(mock.Anything).Return([]*dto.User{{ID: 1, Username: "aaron"}}, nil)
	return &User{svc: svc}
}

func TestGetUserList(t *testing.T) {
	assert := assert.New(t)
	user := mockUser(t)

	resp := test.NewRequest(t).Handler(user.List).Get("/user/list").Exec()
	assert.Equal(200, resp.Code())
	errno, _, _, err := resp.TryDecode()
	assert.Equal(errno, constant.CodeSuccess)
	assert.Nil(err)
}
