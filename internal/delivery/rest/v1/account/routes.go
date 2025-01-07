package account

import (
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/account/verify"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(v1 *gin.RouterGroup) {
	accounts := v1.Group("/accounts")
	verify.Register(accounts)
}
