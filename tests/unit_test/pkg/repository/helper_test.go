package repository_test

import (
	"time"

	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/transaction"
	"github.com/banggok/boillerplate_architecture/internal/pkg/repository"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

func setupTestDB() *gorm.DB {
	dbLogger := logger.New(
		log.New(),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		panic("failed to connect to test database")
	}

	return db
}

func setupTransaction(slaveDB *gorm.DB, masterDB *gorm.DB) *gin.Context {
	read := slaveDB.Session(&gorm.Session{AllowGlobalUpdate: false,
		PrepareStmt: true,
	})
	tx := masterDB.Session(&gorm.Session{
		PrepareStmt: true,
	})

	ctx := &gin.Context{}
	ctx.Set(transaction.DBREAD, read)
	ctx.Set(transaction.DBTRANSACTION, tx)
	return ctx
}

func setupTestModel() (*gorm.DB, *gin.Context, repository.GenericRepository[TestModel, TestModel]) {
	db := setupTestDB()
	db.AutoMigrate(&TestModel{})
	ctx := setupTransaction(db, db)

	testModel := getTestModel()

	db.Create(&testModel[0])
	db.Create(&testModel[1])

	repo := repository.NewGenericRepository(func(entity TestModel) TestModel {
		return entity
	})
	return db, ctx, repo
}

func getTestModel() []TestModel {
	return []TestModel{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
	}
}
