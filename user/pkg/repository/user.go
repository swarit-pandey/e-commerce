package repository

import (
	"context"
	"errors"

	"github.com/swarit-pandey/e-commerce/common/database"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, request *User) (*User, error)
	List(ctx context.Context, request *User) ([]User, error)
	Exists(ctx context.Context, request *User) (bool, error)
	Get(ctx context.Context, request *User) (*User, error)
	GetByUsername(ctx context.Context, request *User) (*User, error)
	Update(ctx context.Context, request *User) (*User, error)
	Delete(ctx context.Context, request *User) (*User, error)
}

type userRepo struct {
	db      *gorm.DB
	handler *database.SQLHandler[User]
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db:      db,
		handler: database.NewSQLHandler[User](db),
	}
}

func (ur *userRepo) Create(ctx context.Context, request *User) (*User, error) {
	err := ur.handler.Create(ctx, request)
	if err != nil {
		return nil, errors.Join(ErrCreate, err)
	}

	return request, nil
}

func (ur *userRepo) List(ctx context.Context, request *User) ([]User, error) {
	resp, err := ur.handler.All(ctx)
	if err != nil || len(resp) == 0 {
		return nil, errors.Join(ErrList, err)
	}

	return resp, nil
}

func (ur *userRepo) Exists(ctx context.Context, request *User) (bool, error) {
	qf := database.QueryFilter{
		WhereMap: map[string]any{
			"username": request.Username,
			"email":    request.Email,
		},
	}

	exists, err := ur.handler.Exists(ctx, qf)
	if err != nil {
		return false, errors.Join(ErrExists, err)
	}

	return exists, nil
}

func (ur *userRepo) Get(ctx context.Context, request *User) (*User, error) {
	qf := database.QueryFilter{
		WhereMap: map[string]any{
			"id": request.ID,
		},
	}

	resp, err := ur.handler.Get(ctx, qf, nil)
	if err != nil {
		return nil, errors.Join(ErrRead, err)
	}

	return resp, nil
}

func (ur *userRepo) GetByUsername(ctx context.Context, request *User) (*User, error) {
	qf := database.QueryFilter{
		WhereMap: map[string]any{
			"username": request.Username,
		},
	}

	resp, err := ur.handler.Get(ctx, qf, nil)
	if err != nil {
		return nil, errors.Join(ErrRead, err)
	}
	return resp, nil
}

func (ur *userRepo) Update(ctx context.Context, request *User) (*User, error) {
	qf := database.QueryFilter{
		WhereMap: map[string]any{
			"id": request.ID,
		},
		UpdateMap: map[string]any{
			"name":  request.Name,
			"email": request.Email,
		},
	}

	err := ur.handler.Update(ctx, qf)
	if err != nil {
		return nil, errors.Join(ErrUpdate, err)
	}

	return request, nil
}

func (ur *userRepo) Delete(ctx context.Context, request *User) (*User, error) {
	qf := database.QueryFilter{
		WhereMap: map[string]any{
			"id": request.ID,
		},
	}

	err := ur.handler.Delete(ctx, qf)
	if err != nil {
		return nil, errors.Join(ErrDelete, err)
	}

	return request, nil
}
