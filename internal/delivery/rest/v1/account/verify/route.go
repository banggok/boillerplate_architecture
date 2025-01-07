package verify

import (
	"github.com/banggok/boillerplate_architecture/internal/services"
	"github.com/gin-gonic/gin"
)

func Register(accountRouterGroup *gin.RouterGroup) {
	usecase := newUsecase(services.ServiceConFig.AccountVerification())
	handler := new(usecase)
	accountRouterGroup.GET("/verify/:token", handler.verify)
}
