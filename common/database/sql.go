package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

var _ Repository[any] = &SQLHandler[any]{}

// SQLHandler wraps a GORM DB connection, providing CRUD operations for a generic type T.
type SQLHandler[T any] struct {
	db *gorm.DB
}

// NewSQLHandler creates a new instance of SQLHandler for the given GORM DB connection.
func NewSQLHandler[T any](db *gorm.DB) *SQLHandler[T] {
	return &SQLHandler[T]{db: db}
}

// Create inserts a new entity of type T into the database.
func (r *SQLHandler[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update modifies entities of type T in the database based on the provided QueryFilter.
func (r *SQLHandler[T]) Update(ctx context.Context, qf QueryFilter) error {
	db := r.db.WithContext(ctx).Model(new(T))
	db = applyQueryFilters(db, qf)
	return db.Updates(qf.UpdateMap).Error
}

// Delete removes entities of type T from the database based on the provided QueryFilter.
func (r *SQLHandler[T]) Delete(ctx context.Context, qf QueryFilter) error {
	db := r.db.WithContext(ctx).Model(new(T))
	db = applyQueryFilters(db, qf)
	return db.Delete(new(T)).Error
}

// Get retrieves a single entity of type T from the database matching the QueryFilter and Projection.
func (r *SQLHandler[T]) Get(ctx context.Context, qf QueryFilter, p *Projection) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx).Model(new(T))
	db = applyQueryFilters(db, qf)
	db = applyProjections(db, p)
	err := db.First(&entity).Error
	return &entity, err
}

// List retrieves multiple entities of type T matching the QueryFilter and Projection.
func (r *SQLHandler[T]) List(ctx context.Context, qf QueryFilter, p *Projection) ([]T, error) {
	var entities []T
	db := r.db.WithContext(ctx).Model(new(T))
	db = applyQueryFilters(db, qf)
	db = applyProjections(db, p)
	err := db.Find(&entities).Error
	return entities, err
}

// Batch retrieves multiple entities of type T by their IDs.
func (r *SQLHandler[T]) Batch(ctx context.Context, columnName string, ids []any) ([]T, error) {
	var entities []T
	db := r.db.WithContext(ctx).Model(new(T))

	err := db.Where(fmt.Sprintf("%s IN (?)", columnName), ids).Find(&entities).Error
	return entities, err
}

// Count returns the count of entities of type T matching the QueryFilter.
func (r *SQLHandler[T]) Count(ctx context.Context, qf QueryFilter) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(new(T))
	db = applyQueryFilters(db, qf)
	err := db.Count(&count).Error
	return count, err
}

// Exists checks if any entity of type T matches the QueryFilter.
func (r *SQLHandler[T]) Exists(ctx context.Context, qf QueryFilter) (bool, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(new(T))
	db = applyQueryFilters(db, qf)
	err := db.Limit(1).Count(&count).Error
	return count > 0, err
}

// All retrieves all entities of type T from the database.
func (r *SQLHandler[T]) All(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}
