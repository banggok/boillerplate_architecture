package repository

import (
	"fmt"

	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/transaction"
	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page int
	Size int
}

type retrievalRepository[E any, M any] interface {
	GetOne(ctx *gin.Context) (E, error)
	GetAll(ctx *gin.Context) (*[]E, error)
	GetAllWithPagination(ctx *gin.Context, pagination Pagination) (*[]E, int64, error)
}

// GetOne retrieves a single entity based on the provided condition.
func (r *genericRepositoryImpl[E, M]) GetOne(ctx *gin.Context) (E, error) {
	// Get database transaction from context
	var entity E
	tx, err := transaction.GetTransaction(ctx, false)
	if err != nil {
		return entity, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Retrieve the model using the condition
	var model M
	tx = r.getQuery(tx)
	if err := tx.First(&model).Error; err != nil {
		return entity, custom_errors.New(err, r.notFoundError(model), "not found")
	}

	// Convert model to entity using the model's ToEntity method
	return r.toEntity(model)
}

// GetAll retrieves all entities based on the provided conditions.
func (r *genericRepositoryImpl[E, M]) GetAll(ctx *gin.Context) (*[]E, error) {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Retrieve all models matching the conditions
	var models []M
	tx = r.getQuery(tx)
	if err := tx.Find(&models).Error; err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to retrieve models from database")
	}

	// Convert models to entities
	var entities []E
	for _, model := range models {
		entity, err := r.toEntity(model)
		if err != nil {
			return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to convert model to entity")
		}
		entities = append(entities, entity)
	}

	return &entities, nil
}

// GetAllWithPagination retrieves all entities with pagination support.
func (r *genericRepositoryImpl[E, M]) GetAllWithPagination(ctx *gin.Context, pagination Pagination) (*[]E, int64, error) {
	// Validate pagination parameters
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Size <= 0 {
		pagination.Size = 10 // Default page size
	}

	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx, false)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Retrieve the total count for the given conditions
	var total int64
	tx = r.getQuery(tx)
	if err := tx.Model(new(M)).Count(&total).Error; err != nil {
		return nil, 0, custom_errors.New(err, custom_errors.InternalServerError, "failed to count rows with given conditions")
	}

	// Retrieve the models with pagination
	var models []M
	offset := (pagination.Page - 1) * pagination.Size
	if err := tx.Limit(pagination.Size).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, custom_errors.New(err, custom_errors.InternalServerError, "failed to retrieve models from database")
	}

	// Convert models to entities
	var entities []E
	for _, model := range models {
		entity, err := r.toEntity(model)
		if err != nil {
			return nil, 0, custom_errors.New(err, custom_errors.InternalServerError, "failed to convert model to entity")
		}
		entities = append(entities, entity)
	}

	return &entities, total, nil
}
