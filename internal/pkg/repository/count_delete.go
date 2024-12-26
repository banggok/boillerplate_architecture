package repository

import (
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/pkg/middleware/transaction"
	"fmt"

	"github.com/gin-gonic/gin"
)

type countDeleteRepository[E any] interface {
	Count(ctx *gin.Context, conditions Condition) (*int64, error)
	Delete(ctx *gin.Context, conditions Condition) error
}

// Count returns the number of rows that match the given condition.
func (r *genericRepositoryImpl[E, M]) Count(ctx *gin.Context, conditions Condition) (*int64, error) {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Count the rows matching the condition
	var count int64
	if err := tx.Model(new(M)).Where(conditions.Query, conditions.Args...).Count(&count).Error; err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to count rows with given condition")
	}

	return &count, nil
}

// Delete removes records based on the provided conditions.
func (r *genericRepositoryImpl[E, M]) Delete(ctx *gin.Context, conditions Condition) error {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Perform the delete operation
	if err := tx.Where(conditions.Query, conditions.Args...).Delete(new(M)).Error; err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to delete records from database")
	}

	return nil
}
