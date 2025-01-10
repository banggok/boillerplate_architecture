package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteWithTransaction(t *testing.T) {
	db, ctx, repo := setupTestModel()

	err := repo.Where("name LIKE ?", "Item%").Delete(ctx)

	require.NoError(t, err)

	// Verify deletion
	var count int64
	db.Model(&TestModel{}).Count(&count)
	assert.Equal(t, int64(0), count)
}
