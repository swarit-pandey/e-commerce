package database

import "context"

// Repository represents a generic interface for CRUD operations on a database entity.
// It defines methods to Create, Update, Delete, and Retrieve entities from the database.
type Repository[T any] interface {
	// Create inserts a new entity into the database.
	// The operation is performed within the context's scope and may return an error.
	Create(ctx context.Context, entity *T) error

	// Update modifies an existing entity based on set of conditions defined as query filter.
	// It returns an error if the update operation fails.
	Update(ctx context.Context, qf QueryFilter) error

	// Delete removes an entity from the database based on where conditions.
	// It returns an error if the deletion fails.
	Delete(ctx context.Context, qf QueryFilter) error

	// Get retrieves a single entity from the database that matches the provided query filters and projection.
	// It returns the found entity or an error if the retrieval fails.
	Get(ctx context.Context, qf QueryFilter, p *Projection) (*T, error)

	// List retrieves multiple entities from the database that match the provided query filters and projection.
	// It returns a slice of entities or an error if the retrieval fails.
	List(ctx context.Context, qf QueryFilter, p *Projection) ([]T, error)

	// Batch retrieves a batch of entities from the database based on the provided column name and IDs.
	// It's upto clients to correctly define the column which has IDs (such as the column that has primary keys).
	// It returns a slice of entities or an error if the retrieval fails.
	Batch(ctx context.Context, columnName string, ids []any) ([]T, error)

	// Count returns the count of entities that match the provided query filters.
	// It returns the count or an error if the operation fails.
	Count(ctx context.Context, qf QueryFilter) (int64, error)

	// Exists checks if there is at least one entity in the database that matches the given conditions.
	// It returns a boolean indicating existence and an error if the check fails.
	Exists(ctx context.Context, qf QueryFilter) (bool, error)

	// All returns all the entities from the database without any filters.
	// It returns slice of entities or an error if the retrieval fails.
	All(ctx context.Context) ([]T, error)
}
