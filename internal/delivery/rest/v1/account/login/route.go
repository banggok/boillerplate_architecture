package login

import (
	"github.com/banggok/boillerplate_architecture/internal/services"
	"github.com/gin-gonic/gin"
)

func Register(accounts *gin.RouterGroup) {
	uc := newUsecase(services.ServiceConFig.Account())
	h := new(uc)
	accounts.POST("/login", h.login)
}
