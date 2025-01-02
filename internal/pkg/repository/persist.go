package repository

import (
	"fmt"

	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/transaction"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type persistRepository[E any] interface {
	Persist(ctx *gin.Context, entity *E) error
}

// Insert inserts the given entity into the database.
func (r *genericRepositoryImpl[E, M]) Persist(ctx *gin.Context, entity *E) error {
	// Convert entity to model
	model := r.toModel(*entity)

	// Get database transaction from context
	tx, err := transaction.GetTransaction(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get database transaction: %w", err)
	}

	fmt.Printf("insert: %v\n", model)
	// Save model to database
	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&model).Error; err != nil {
		fmt.Printf("error_persist: %v\n", err)
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to save model to database")
	}

	// Convert model back to entity using the model's ToEntity method
	updatedEntity, err := r.toEntity(model)

	if err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to convert model to entity")
	}

	*entity = updatedEntity

	return nil
}
