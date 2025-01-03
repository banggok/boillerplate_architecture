package register

import (
	"net/http"

	"github.com/banggok/boillerplate_architecture/internal/config/validator"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/gin-gonic/gin"
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
	var request Request

	// Parse and validate the request payload
	if err := h.parseAndValidateRequest(c, &request); err != nil {
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

func (h *handlerImpl) parseAndValidateRequest(c *gin.Context, request *Request) error {
	// Parse the incoming JSON request
	if err := c.ShouldBindJSON(request); err != nil {
		c.Error(custom_errors.New(
			err,
			custom_errors.TenantBadRequest,
			"invalid tenant register request payload"))
		return err
	}

	// Validate the request using the validator
	if err := validator.Validate.Struct(request); err != nil {
		validationErrors := request.customValidationMessages(err)
		c.Error(custom_errors.New(
			err,
			custom_errors.TenantUnprocessEntity,
			"failed to validate request payload",
			validationErrors))
		return err
	}

	return nil
}
