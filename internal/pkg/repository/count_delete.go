package repository

import (
	"fmt"

	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/transaction"
	"github.com/gin-gonic/gin"
)

type countDeleteRepository[E any] interface {
	Count(ctx *gin.Context) (*int64, error)
	Delete(ctx *gin.Context) error
}

// Count returns the number of rows that match the given condition.
func (r *genericRepositoryImpl[E, M]) Count(ctx *gin.Context) (*int64, error) {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Count the rows matching the condition
	var count int64
	tx = r.getWhere(tx)
	if err := tx.Model(new(M)).Count(&count).Error; err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to count rows with given condition")
	}

	return &count, nil
}

// Delete removes records based on the provided conditions.
func (r *genericRepositoryImpl[E, M]) Delete(ctx *gin.Context) error {
	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Perform the delete operation
	tx = r.getWhere(tx)
	if err := tx.Delete(new(M)).Error; err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to delete records from database")
	}

	return nil
}
