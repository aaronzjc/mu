package store

import (
	"context"
	"errors"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/domain/repo"
	"gorm.io/gorm"
)

type UserRepoImpl struct {
	BaseRepoImpl
}

var _ repo.UserRepo = &UserRepoImpl{}

func (r *UserRepoImpl) GetAll(ctx context.Context) ([]model.User, error) {
	users := []model.User{}
	st := r.db.Find(&users)
	if st.Error != nil {
		return nil, errors.New("get users err")
	}
	return users, nil
}

func (r *UserRepoImpl) GetUser(ctx context.Context, query *dto.Query) (model.User, error) {
	var user model.User
	err := r.prepare(query).First(&user).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return model.User{}, nil
		}
		return model.User{}, errors.New("fetch user info failed")
	}

	return user, nil
}

func (r *UserRepoImpl) GetUsers(ctx context.Context, query *dto.Query) ([]model.User, error) {
	var users []model.User
	err := r.prepare(query).Find(&users).Error
	if err != nil {
		return nil, errors.New("fetch user info failed")
	}

	return users, nil
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, user model.User) error {
	var exist model.User
	var err error
	exist, err = r.GetUser(ctx, &dto.Query{
		Query: "`username` = ?",
		Args:  []interface{}{user.Username},
	})
	if err != nil {
		return errors.New("find user err")
	}
	if exist.ID > 0 {
		return errors.New("user exist")
	}
	if err = r.create(&user); err != nil {
		return errors.New("create user err")
	}
	return nil
}

func (r *UserRepoImpl) Update(ctx context.Context, user model.User, data map[string]interface{}) error {
	return r.db.Model(&user).Updates(data).Error
}

func NewUserRepo() *UserRepoImpl {
	base, _ := NewBaseImpl()
	return &UserRepoImpl{
		base,
	}
}
