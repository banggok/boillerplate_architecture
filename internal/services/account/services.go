package account

import (
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/pkg/repository"
	"github.com/gin-gonic/gin"
)

type Service interface {
	GetAccountVerifiedByEmail(ctx *gin.Context, email string) (entity.Account, error)
}

type serviceImpl struct {
	repo repository.GenericRepository[entity.Account, model.Account]
}

// GetAccountVerifiedByEmail implements Service.
func (s *serviceImpl) GetAccountVerifiedByEmail(ctx *gin.Context, email string) (entity.Account, error) {
	return s.repo.Where("email = ?", email).Preload("AccountVerifications").GetOne(ctx)
}

func New(repo repository.GenericRepository[entity.Account, model.Account]) Service {
	return &serviceImpl{
		repo: repo,
	}
}
