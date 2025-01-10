package repository

import (
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"gorm.io/gorm"
)

type queryArgs struct {
	Query string
	Args  []interface{}
}

// GenericRepository defines the interface for CRUD operations.
type GenericRepository[E any, M any] interface {
	persistRepository[E]
	retrievalRepository[E, M]
	countDeleteRepository[E]

	// Where appends a WHERE clause to the repository's query builder.
	// This allows filtering results based on specific conditions.
	//
	// Parameters:
	//   - query: A string representing the condition for the WHERE clause
	//     (e.g., "age > ? AND status = ?").
	//   - args: Optional arguments to bind to the placeholders in the query string.
	//
	// Returns:
	//   - GenericRepository[E, M]: The repository instance with the new WHERE condition appended,
	//     enabling method chaining.
	//
	// Example Usage:
	//
	//	repo.Where("age > ?", 30).Where("status = ?", "active")
	//
	// Notes:
	//   - This method is typically used to refine database queries by specifying
	//     conditions that must be met by the resulting rows.
	Where(query string, args ...interface{}) GenericRepository[E, M]

	// Joins appends an INNER JOIN clause to the repository's query builder. For LEFT JOIN use Preload().
	// It allows specifying the join condition as a query string and any associated arguments.
	//
	// Parameters:
	//   - query: A string representing the INNER JOIN condition (e.g., "JOIN users ON users.id = orders.user_id").
	//   - args: Optional arguments to bind to the query string, typically for placeholders (e.g., WHERE conditions).
	//
	// Returns:
	//   - GenericRepository[E, M]: The repository instance with the new join condition appended, allowing method chaining.
	//
	// Example Usage:
	//
	//	repo.Joins("Account")
	//	repo.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org")
	//	repo.Joins("Account").Where("user_id = users.id AND name = ?", "someName")
	// Joins(query string, args ...interface{}) GenericRepository[E, M]

	// Preload appends a preload operation to the repository's query builder.
	// This allows eager loading of related entities or associations with optional conditions.
	//
	// Parameters:
	//   - query: A string representing the association to preload
	//     (e.g., "Orders", "Orders.Items").
	//   - args: Optional arguments for conditions or additional filtering
	//     applied to the preload operation.
	//
	// Returns:
	//   - GenericRepository[E, M]: The repository instance with the new preload condition appended,
	//     allowing method chaining.
	//
	// Example Usage:
	//
	//	// Basic preloading of an association
	//	repo.Preload("Orders")
	//
	//	// Preloading with a condition to filter the associated records
	//	repo.Preload("Orders", "status = ?", "completed")
	//
	// Notes:
	//   - Preloading is typically used for LEFT JOIN-like behavior, ensuring related entities
	//     are loaded alongside the main entity in a single query or a sequence of queries.
	//   - Conditional preloading is useful for reducing the size of the loaded dataset and focusing
	//     only on relevant related data.
	Preload(query string, args ...interface{}) GenericRepository[E, M]
}

// genericRepositoryImpl implements the GenericRepository interface.
type genericRepositoryImpl[E any, M any] struct {
	// mu      sync.Mutex // Mutex for synchronizing access
	toModel func(entity E) M
	wheres  []queryArgs
	joins   []queryArgs
	preload []queryArgs
}

func (r *genericRepositoryImpl[E, M]) Where(query string, args ...interface{}) GenericRepository[E, M] {
	// r.mu.Lock()
	// defer r.mu.Unlock()
	r.wheres = append(r.wheres, queryArgs{
		Query: query,
		Args:  args,
	})

	return r
}

// func (r *genericRepositoryImpl[E, M]) Joins(query string, args ...interface{}) GenericRepository[E, M] {
// r.mu.Lock()
// defer r.mu.Unlock()

// 	r.joins = append(r.joins, queryArgs{
// 		Query: query,
// 		Args:  args,
// 	})

// 	return r
// }

func (r *genericRepositoryImpl[E, M]) Preload(query string, args ...interface{}) GenericRepository[E, M] {
	// r.mu.Lock()
	// defer r.mu.Unlock()

	r.preload = append(r.preload, queryArgs{
		Query: query,
		Args:  args,
	})

	return r
}

func (r *genericRepositoryImpl[E, M]) getQuery(tx *gorm.DB) *gorm.DB {
	// r.mu.Lock()
	// defer r.mu.Unlock()
	for _, where := range r.wheres {
		tx = tx.Where(where.Query, where.Args...)
	}

	for _, join := range r.joins {
		tx = tx.Joins(join.Query, join.Args...)
	}

	for _, preload := range r.preload {
		tx = tx.Preload(preload.Query, preload.Args...)
	}
	return tx
}

func (r *genericRepositoryImpl[E, M]) toEntity(model M) (E, error) {
	var entity E
	if toEntityMethod, ok := any(&model).(interface{ ToEntity() (E, error) }); ok {
		entity, err := toEntityMethod.ToEntity()
		if err != nil {
			return entity, custom_errors.New(err, custom_errors.InternalServerError, "failed to convert model to entity")
		}
		return entity, nil
	}
	return entity, custom_errors.New(nil, custom_errors.InternalServerError, "model does not implement ToEntity method")
}

func (r *genericRepositoryImpl[E, M]) notFoundError(model M) custom_errors.ErrorCode {
	if toEntityMethod, ok := any(&model).(interface {
		NotFoundError() custom_errors.ErrorCode
	}); ok {
		return toEntityMethod.NotFoundError()
	}
	return custom_errors.InternalServerError
}

// NewGenericRepository creates a new instance of GenericRepository.
func NewGenericRepository[E any, M any](toModel func(E) M) GenericRepository[E, M] {
	return &genericRepositoryImpl[E, M]{
		toModel: toModel,
	}
}
