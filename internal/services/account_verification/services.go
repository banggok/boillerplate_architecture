package accountverification

import (
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/repository"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Verify(*gin.Context, *entity.AccountVerification) error
	GetByTokenVerification(ctx *gin.Context, token string) (entity.AccountVerification, error)
}

type serviceImpl struct {
	repo repository.GenericRepository[entity.AccountVerification, model.AccountVerification]
}

// GetByTokenVerification implements Service.
func (s *serviceImpl) GetByTokenVerification(ctx *gin.Context, token string) (entity.AccountVerification, error) {
	return s.repo.Where("token = ?", token).Preload("Account").GetOne(ctx)
}

// VerifySuccess implements Service.
func (s *serviceImpl) Verify(ctx *gin.Context, accountVerificationEntity *entity.AccountVerification) error {
	if accountVerificationEntity == nil {
		return custom_errors.New(nil, custom_errors.AccountUnprocessEntity, "accountVerificationEntity was empty")
	}

	(*accountVerificationEntity).VerifiedSuccess()

	return s.repo.Persist(ctx, accountVerificationEntity)
}

func NewService(accountVerificationRepo repository.GenericRepository[entity.AccountVerification, model.AccountVerification]) Service {
	return &serviceImpl{
		repo: accountVerificationRepo,
	}
}
