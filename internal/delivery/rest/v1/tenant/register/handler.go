package register

import (
	"net/http"

	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// handler represents the HTTP handler for tenant-related operations
type handler interface {
	register(*gin.Context)
}

type handlerImpl struct {
	usecase usecase
}

// NewTenantRegisterHandler creates a newHandler TenantHandler instance with a validator
func newHandler(
	uc usecase,
) handler {
	return &handlerImpl{
		usecase: uc,
	}
}

// register handles tenant registration
// @Summary register Tenant
// @Description register a new tenant in the system
// @Tags Tenants
// @Accept json
// @Produce json
// @Param body body Request true "Tenant registration request body"
// @Success 201 {object} Response
// @Failure 400 {object} custom_errors.CustomError "Invalid request payload"
// @Failure 422 {object} custom_errors.CustomError "Failed to register tenant"
// @Router /api/v1/tenants/ [post]
func (h *handlerImpl) register(c *gin.Context) {
	request := Request{}

	// Parse and validate the request payload
	if err := request.ParseAndValidateRequest(c, binding.JSON); err != nil {
		return
	}

	// Execute the usecase
	tenant, err := h.usecase.execute(c, request)
	if err != nil {
		c.Error(custom_errors.New(err, custom_errors.TenantUnprocessEntity, "failed to register the tenant."))
		return
	}

	// Respond with the created tenant
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Tenant registered successfully",
		"data":    transform(tenant),
	})
}
