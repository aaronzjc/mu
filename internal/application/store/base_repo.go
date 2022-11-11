package store

import (
	"errors"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/infra/db"
	"gorm.io/gorm"
)

type BaseRepoImpl struct {
	db *gorm.DB
}

func (r *BaseRepoImpl) prepare(q *dto.Query) *gorm.DB {
	clone := r.db
	if q == nil {
		return clone
	}
	if q.Order != "" {
		clone = clone.Order(q.Order)
	}

	if q.Query != "" {
		clone = clone.Where(q.Query, q.Args...)
	}

	if q.Limit > 0 {
		clone = clone.Limit(q.Limit)
	}
	return clone
}

func (r *BaseRepoImpl) create(m interface{}) error {
	return r.db.Create(m).Error
}

func NewBaseImpl() (BaseRepoImpl, error) {
	mu, ok := db.Get(model.DB_MU)
	if !ok {
		return BaseRepoImpl{}, errors.New("db not connected")
	}
	return BaseRepoImpl{db: mu}, nil
}
