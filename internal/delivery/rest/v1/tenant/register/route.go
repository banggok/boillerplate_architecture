package register

import (
	"github.com/banggok/boillerplate_architecture/internal/services"
	"github.com/gin-gonic/gin"
)

func Register(
	tenantV1RouterGroup *gin.RouterGroup,
) {
	uc := newUsecase(services.ServiceConFig.Tenant(), services.ServiceConFig.Email())
	h := newHandler(uc)
	tenantV1RouterGroup.POST("", h.register)
}
