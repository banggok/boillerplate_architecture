package repository

import (
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"gorm.io/gorm"
)

type condition struct {
	Query string
	Args  []interface{}
}

// GenericRepository defines the interface for CRUD operations.
type GenericRepository[E any, M any] interface {
	persistRepository[E]
	retrievalRepository[E, M]
	countDeleteRepository[E]
	Where(query string, args ...interface{}) GenericRepository[E, M]
}

// genericRepositoryImpl implements the GenericRepository interface.
type genericRepositoryImpl[E any, M any] struct {
	toModel    func(entity E) M
	conditions []condition
}

func (r *genericRepositoryImpl[E, M]) Where(query string, args ...interface{}) GenericRepository[E, M] {
	r.conditions = append(r.conditions, condition{
		Query: query,
		Args:  args,
	})

	return r
}

func (r *genericRepositoryImpl[E, M]) getWhere(tx *gorm.DB) *gorm.DB {
	for _, condition := range r.conditions {
		tx = tx.Where(condition.Query, condition.Args...)
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
