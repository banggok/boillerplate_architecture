package account

import (
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/pkg/repository"
	"github.com/gin-gonic/gin"
)

type Service interface {
	GetById(ctx *gin.Context, id uint) (entity.Account, error)
}

type serviceImpl struct {
	repo repository.GenericRepository[entity.Account, model.Account]
}

// GetById implements Service.
func (s *serviceImpl) GetById(ctx *gin.Context, id uint) (entity.Account, error) {
	return s.repo.Where("id = ?", id).GetOne(ctx)
}

func New(repo repository.GenericRepository[entity.Account, model.Account]) Service {
	return &serviceImpl{
		repo: repo,
	}
}
