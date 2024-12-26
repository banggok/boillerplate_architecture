package repository

import (
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/pkg/middleware/transaction"
	"fmt"

	"github.com/gin-gonic/gin"
)

type retrievalRepository[E any, M any] interface {
	GetOne(ctx *gin.Context, conditions Condition) (*E, error)
	GetAll(ctx *gin.Context, conditions Condition) (*[]E, error)
	GetAllWithPagination(ctx *gin.Context, conditions Condition, page int, pageSize int) (*[]E, int64, error)
}

// GetOne retrieves a single entity based on the provided condition.
func (r *genericRepositoryImpl[E, M]) GetOne(ctx *gin.Context, conditions Condition) (*E, error) {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx, false)
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

// GetAll retrieves all entities based on the provided conditions.
func (r *genericRepositoryImpl[E, M]) GetAll(ctx *gin.Context, conditions Condition) (*[]E, error) {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx, false)
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
	tx, err := transaction.GetTransaction(ctx, false)
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
