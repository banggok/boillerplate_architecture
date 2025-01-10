package login

import "github.com/gin-gonic/gin"

func Register(accounts *gin.RouterGroup) {
	uc := newUsecase()
	h := new(uc)
	accounts.POST("/login", h.login)
}
