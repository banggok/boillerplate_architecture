package login

import (
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/gin-gonic/gin"
)

type handler interface {
	login(*gin.Context)
}

type handlerImpl struct {
	usecase usecase
}

// login implements handler.
func (h *handlerImpl) login(c *gin.Context) {
	request := Request{}

	// Parse and validate the request payload
	if err := request.ParseAndValidateRequest(c); err != nil {
		return
	}

	// Execute the usecase
	_, err := h.usecase.execute(c, request)
	if err != nil {
		c.Error(custom_errors.New(err, custom_errors.AccountUnprocessEntity, "failed to login account."))
		return
	}

	// Respond with the created tenant
	// c.JSON(http.StatusCreated, gin.H{
	// 	"status":  "success",
	// 	"message": "Tenant registered successfully",
	// 	"data":    transform(tenant),
	// })
}

func new(usecase usecase) handler {
	return &handlerImpl{
		usecase: usecase,
	}
}
