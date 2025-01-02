package repository_test

import (
	"appointment_management_system/internal/pkg/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteWithTransaction(t *testing.T) {
	db, ctx, repo := setupTestModel()

	conditions := repository.Condition{
		Query: "name LIKE ?",
		Args:  []interface{}{"Item%"},
	}

	err := repo.Delete(ctx, conditions)

	assert.NoError(t, err)

	// Verify deletion
	var count int64
	db.Model(&TestModel{}).Count(&count)
	assert.Equal(t, int64(0), count)
}