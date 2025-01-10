package login

import (
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/request"
	"github.com/gin-gonic/gin"
)

type usecase interface {
	execute(ctx *gin.Context, iRequest request.IRequest) (entity.Account, error)
}

type usecaseImpl struct{}

// execute implements usecase.
func (u *usecaseImpl) execute(ctx *gin.Context, iRequest request.IRequest) (entity.Account, error) {
	panic("unimplemented")
}

func newUsecase() usecase {
	return &usecaseImpl{}
}
