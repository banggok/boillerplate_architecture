package repository

import (
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/pkg/middleware/transaction"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Condition struct {
	Query string
	Args  []interface{}
}

type persistRepository[E any] interface {
	Persist(ctx *gin.Context, entity *E) error
}

type countDeleteRepository[E any] interface {
	Count(ctx *gin.Context, conditions Condition) (int64, error)
	Delete(ctx *gin.Context, conditions Condition) error
}

type retrievalRepository[E any, M any] interface {
	GetOne(ctx *gin.Context, conditions Condition) (*E, error)
	GetAll(ctx *gin.Context, conditions Condition) (*[]E, error)
	GetAllWithPagination(ctx *gin.Context, conditions Condition, page int, pageSize int) (*[]E, int64, error)
}

// GenericRepository defines the interface for CRUD operations.
type GenericRepository[E any, M any] interface {
	persistRepository[E]
	retrievalRepository[E, M]
	countDeleteRepository[E]
}

// genericRepositoryImpl implements the GenericRepository interface.
type genericRepositoryImpl[E any, M any] struct {
	toModel func(entity E) M
}

// Insert inserts the given entity into the database.
func (r *genericRepositoryImpl[E, M]) Persist(ctx *gin.Context, entity *E) error {
	// Convert entity to model
	model := r.toModel(*entity)

	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Save model to database
	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&model).Error; err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to save model to database")
	}

	// Convert model back to entity using the model's ToEntity method
	updatedEntity, err := r.toEntity(model)

	if err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to convert model to entity")
	}

	*entity = *updatedEntity

	return nil
}

// GetOne retrieves a single entity based on the provided condition.
func (r *genericRepositoryImpl[E, M]) GetOne(ctx *gin.Context, conditions Condition) (*E, error) {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Retrieve the model using the condition
	var model M
	if err := tx.Where(conditions.Query, conditions.Args...).First(&model).Error; err != nil {
		return nil, custom_errors.New(err, r.notFoundError(model), "failed to find model with given condition")
	}

	// Convert model to entity using the model's ToEntity method
	return r.toEntity(model)
}

// Count returns the number of rows that match the given condition.
func (r *genericRepositoryImpl[E, M]) Count(ctx *gin.Context, conditions Condition) (int64, error) {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Count the rows matching the condition
	var count int64
	if err := tx.Model(new(M)).Where(conditions.Query, conditions.Args...).Count(&count).Error; err != nil {
		return 0, custom_errors.New(err, custom_errors.InternalServerError, "failed to count rows with given condition")
	}

	return count, nil
}

// GetAll retrieves all entities based on the provided conditions.
func (r *genericRepositoryImpl[E, M]) GetAll(ctx *gin.Context, conditions Condition) (*[]E, error) {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Retrieve all models matching the conditions
	var models []M
	if err := tx.Where(conditions.Query, conditions.Args...).Find(&models).Error; err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to retrieve models from database")
	}

	// Convert models to entities
	var entities []E
	for _, model := range models {
		entity, err := r.toEntity(model)
		if err != nil {
			return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to convert model to entity")
		}
		entities = append(entities, *entity)
	}

	return &entities, nil
}

// Delete removes records based on the provided conditions.
func (r *genericRepositoryImpl[E, M]) Delete(ctx *gin.Context, conditions Condition) error {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Perform the delete operation
	if err := tx.Model(new(M)).Where(conditions.Query, conditions.Args...).Delete(new(M)).Error; err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to delete records from database")
	}

	return nil
}

// GetAllWithPagination retrieves all entities with pagination support.
func (r *genericRepositoryImpl[E, M]) GetAllWithPagination(ctx *gin.Context, conditions Condition, page int, pageSize int) (*[]E, int64, error) {
	// Validate pagination parameters
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Retrieve the total count for the given conditions
	var total int64
	if err := tx.Model(new(M)).Where(conditions.Query, conditions.Args...).Count(&total).Error; err != nil {
		return nil, 0, custom_errors.New(err, custom_errors.InternalServerError, "failed to count rows with given conditions")
	}

	// Retrieve the models with pagination
	var models []M
	offset := (page - 1) * pageSize
	if err := tx.Where(conditions.Query, conditions.Args...).Limit(pageSize).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, custom_errors.New(err, custom_errors.InternalServerError, "failed to retrieve models from database")
	}

	// Convert models to entities
	var entities []E
	for _, model := range models {
		entity, err := r.toEntity(model)
		if err != nil {
			return nil, 0, custom_errors.New(err, custom_errors.InternalServerError, "failed to convert model to entity")
		}
		entities = append(entities, *entity)
	}

	return &entities, total, nil
}

func (r *genericRepositoryImpl[E, M]) toEntity(model M) (*E, error) {
	if toEntityMethod, ok := any(&model).(interface{ ToEntity() (*E, error) }); ok {
		entity, err := toEntityMethod.ToEntity()
		if err != nil {
			return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to convert model to entity")
		}
		return entity, nil
	}
	return nil, custom_errors.New(nil, custom_errors.InternalServerError, "model does not implement ToEntity method")
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
