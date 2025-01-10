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

// verify handles account verification
// @Summary Verify Account
// @Description Verifies an account using a token provided in the URL path
// @Tags Accounts
// @Accept json
// @Produce json
// @Param token path string true "Verification token"
// @Success 200 {object} Response
// @Failure 400 {object} custom_errors.CustomError "Invalid request payload"
// @Failure 422 {object} custom_errors.CustomError "Failed to verify account"
// @Router /api/v1/accounts/verify/{token} [get]
func (h *handlerImpl) verify(c *gin.Context) {
	request := newRequest()

	// Parse and validate the request payload
	request.ParseAndValidateRequest(c)

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
