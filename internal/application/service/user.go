package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/domain/repo"
	"github.com/aaronzjc/mu/internal/util"
	"github.com/aaronzjc/mu/pkg/oauth"
)

type UserService interface {
	GetUserList(context.Context) ([]*dto.User, error)
	GetUser(context.Context, *dto.Query) (*dto.User, error)
	Auth(context.Context, string, oauth.User) (string, error)
	VerifyToken(context.Context, string, string) bool
}

type UserServiceImpl struct {
	repo repo.UserRepo
}

func (s *UserServiceImpl) GetUserList(ctx context.Context) ([]*dto.User, error) {
	userModels, err := s.repo.GetUsers(ctx, nil)
	if err != nil {
		return nil, err
	}
	var users []*dto.User
	for _, v := range userModels {
		users = append(users, (&dto.User{}).FillByModel(v))
	}
	return users, nil
}

func (s *UserServiceImpl) GetUser(ctx context.Context, q *dto.Query) (*dto.User, error) {
	user, err := s.repo.GetUser(ctx, q)
	if err != nil {
		return nil, err
	}
	return (&dto.User{}).FillByModel(user), nil
}

func (s *UserServiceImpl) Auth(ctx context.Context, t string, ou oauth.User) (string, error) {
	exist, _ := s.repo.GetUser(ctx, &dto.Query{
		Query: "`username` = ? AND `auth_type` = ?",
		Args:  []interface{}{ou.Username, t},
	})
	conf := config.Get()
	token := util.GenerateToken(ou.Username, conf.Salt)
	expireAt := time.Now().Add(time.Hour * 24 * 30).Unix()
	if exist.ID > 0 {
		if err := s.repo.Update(ctx, exist, map[string]interface{}{
			"token":     token,
			"expire_at": expireAt,
			"auth_time": time.Now(),
		}); err != nil {
			return "", errors.New("auth update err")
		}
	} else {
		user := model.User{
			Username: ou.Username,
			Nickname: ou.Nickname,
			Avatar:   ou.Avatar,
			AuthType: t,
			AuthTime: time.Now(),
			Token:    token,
			ExpireAt: expireAt,
		}
		if err := s.repo.CreateUser(ctx, user); err != nil {
			return "", errors.New("auth create err")
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s;%s", ou.Username, token))), nil
}

func (s *UserServiceImpl) VerifyToken(ctx context.Context, username string, token string) bool {
	login, _ := s.repo.GetUser(ctx, &dto.Query{
		Query: "`username` = ? AND `token` = ?",
		Args:  []interface{}{username, token},
	})
	if login.ID < 0 {
		return false
	}
	if login.ExpireAt <= time.Now().Unix() {
		return false
	}
	return true
}

func NewUserService(repo repo.UserRepo) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}
