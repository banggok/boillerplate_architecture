package login

import (
	"net/http"

	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/authentication"
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
	request := newRequest(c)

	// Parse and validate the request payload
	if err := request.ParseAndValidateRequest(); err != nil {
		return
	}

	// Execute the usecase
	account, err := h.usecase.execute(c, request)
	if err != nil {
		c.Error(custom_errors.New(err, custom_errors.AccountUnprocessEntity, "failed to login account."))
		return
	}

	accessToken, refreshToken, err := authentication.GenerateTokens(account.ID())
	if err != nil {
		c.Error(custom_errors.New(err, custom_errors.InternalServerError, "failed to generate token"))
		return
	}

	data, err := transform(account, accessToken, refreshToken)
	if err != nil || data == nil {
		c.Error(custom_errors.New(err, custom_errors.InternalServerError, "failed to transform data"))
	}

	// Respond with the created tenant
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Account logged in successfully",
		"data":    data,
	})
}

func new(usecase usecase) handler {
	return &handlerImpl{
		usecase: usecase,
	}
}
