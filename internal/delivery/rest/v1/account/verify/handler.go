package verify

import (
	"net/http"

	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/gin-gonic/gin"
)

type handler interface {
	verify(*gin.Context)
}

type handlerImpl struct {
	usecase usecase
}

// verify implements handler.
func (h *handlerImpl) verify(c *gin.Context) {
	request := Request{}

	// Parse and validate the request payload
	request.ParseAndValidateRequest(c, nil)

	// Execute the usecase
	account, err := h.usecase.execute(c, request)
	if err != nil {
		c.Error(custom_errors.New(err, custom_errors.AccountUnprocessEntity, "failed to verify account"))
		return
	}

	// Respond with the created tenant
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Account verified successfully",
		"data":    transform(account),
	})
}

func new(usecase usecase) handler {
	return &handlerImpl{
		usecase: usecase,
	}
}
