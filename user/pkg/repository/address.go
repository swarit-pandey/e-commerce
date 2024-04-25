package repository

import (
	"context"
	"errors"

	"github.com/swarit-pandey/e-commerce/common/database"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

type AddressRepo interface {
	Create(ctx context.Context, request *UserAddress) (*UserAddress, error)
	Get(ctx context.Context, request *UserAddress) (*UserAddress, error)
	Update(ctx context.Context, request *UserAddress) (*UserAddress, error)
	Delete(ctx context.Context, request *UserAddress) (*UserAddress, error)
	Exists(ctx context.Context, request *UserAddress) (bool, error)
}

type addressRepo struct {
	db      *gorm.DB
	handler *database.SQLHandler[UserAddress]
}

func NewAddressRepo(db *gorm.DB) AddressRepo {
	return &addressRepo{
		db:      db,
		handler: database.NewSQLHandler[UserAddress](db),
	}
}

func (ar *addressRepo) Create(ctx context.Context, request *UserAddress) (*UserAddress, error) {
	err := ar.handler.Create(ctx, request)
	if err != nil {
		return nil, errors.Join(ErrCreate, err)
	}

	return request, nil
}

func (ar *addressRepo) Get(ctx context.Context, request *UserAddress) (*UserAddress, error) {
	qf := database.QueryFilter{
		WhereMap: map[string]any{
			"id":      request.ID,
			"user_id": request.UserID,
		},
	}

	resp, err := ar.handler.Get(ctx, qf, nil)
	if err != nil {
		return nil, errors.Join(ErrRead, err)
	}

	return resp, nil
}

func (ar *addressRepo) Update(ctx context.Context, request *UserAddress) (*UserAddress, error) {
	qf := database.QueryFilter{
		WhereMap: map[string]any{
			"id":      request.ID,
			"user_id": request.UserID,
		},
		UpdateMap: map[string]any{
			"address_line1": request.AddressLine1,
			"address_line2": request.AddressLine2,
			"city":          request.City,
			"state":         request.State,
			"country":       request.Country,
			"postal_code":   request.PostalCode,
			"is_primary":    request.IsPrimary,
		},
	}

	err := ar.handler.Update(ctx, qf)
	if err != nil {
		return nil, errors.Join(ErrUpdate, err)
	}

	return request, nil
}

func (ar *addressRepo) Delete(ctx context.Context, request *UserAddress) (*UserAddress, error) {
	qf := database.QueryFilter{
		WhereMap: map[string]any{
			"id":      request.ID,
			"user_id": request.UserID,
		},
	}

	err := ar.handler.Delete(ctx, qf)
	if err != nil {
		return nil, errors.Join(ErrDelete, err)
	}

	return request, nil
}

func (ar *addressRepo) Exists(ctx context.Context, request *UserAddress) (bool, error) {
	qf := database.QueryFilter{
		WhereMap: map[string]any{
			"id": request.ID,
		},
	}

	ok, err := ar.handler.Exists(ctx, qf)
	if err != nil {
		// This should not be false ideally, but lets for now lets return false
		return false, errors.Join(ErrExists, err)
	}

	if !ok {
		klog.Warning("entity might not exist")
		return false, ErrExists
	}

	return true, nil
}
