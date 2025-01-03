package register

import (
	"github.com/banggok/boillerplate_architecture/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(
	tenantV1RouterGroup *gin.RouterGroup,
	serviceConfig services.Config,
) {
	tenantCreateService := serviceConfig.Tenant()
	email := serviceConfig.Email()
	uc := newUsecase(tenantCreateService, email)
	h := newHandler(uc)
	tenantV1RouterGroup.POST("", h.register)
}
