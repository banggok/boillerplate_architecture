package login

import (
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/request"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/password"
	"github.com/banggok/boillerplate_architecture/internal/services/account"
	"github.com/gin-gonic/gin"
)

type usecase interface {
	execute(ctx *gin.Context, iRequest request.IRequest) (entity.Account, error)
}

type usecaseImpl struct {
	service account.Service
}

// execute implements usecase.
func (u *usecaseImpl) execute(ctx *gin.Context, iRequest request.IRequest) (entity.Account, error) {
	request, ok := iRequest.(*Request)
	if !ok {
		return nil, custom_errors.New(nil, custom_errors.InternalServerError, "request invalid")
	}

	account, err := u.service.GetAccountVerifiedByEmail(ctx, request.Email)
	if err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to get account by email")
	}

	match, err := password.VerifyPassword(request.Password, account.Password())
	if err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to match password")
	}

	if !match {
		return nil, custom_errors.New(nil, custom_errors.Unauthorized, "wrong password")
	}

	return account, nil
}

func newUsecase(service account.Service) usecase {
	return &usecaseImpl{
		service: service,
	}
}
