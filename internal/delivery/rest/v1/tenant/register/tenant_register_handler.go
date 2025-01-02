package register

import (
	"net/http"

	"github.com/banggok/boillerplate_architecture/internal/config/validator"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/gin-gonic/gin"
)

// TenantRegisterHandler represents the HTTP handler for tenant-related operations
type TenantRegisterHandler interface {
	Register(*gin.Context)
}

type tenantRegisterHandler struct {
	registerUsecase TenantRegisterUsecase
}

// NewTenantRegisterHandler creates a new TenantHandler instance with a validator
func newTenantRegisterHandler(
	uc TenantRegisterUsecase,
) TenantRegisterHandler {
	return &tenantRegisterHandler{
		registerUsecase: uc,
	}
}

// Register handles tenant registration
// @Summary Register Tenant
// @Description Register a new tenant in the system
// @Tags Tenants
// @Accept json
// @Produce json
// @Param body body RegisterTenantRequest true "Tenant registration request body"
// @Success 201 {object} RegisterTenantResponse
// @Failure 400 {object} custom_errors.CustomError "Invalid request payload"
// @Failure 422 {object} custom_errors.CustomError "Failed to register tenant"
// @Router /api/v1/tenants/ [post]
func (h *tenantRegisterHandler) Register(c *gin.Context) {
	var request RegisterTenantRequest

	// Parse and validate the request payload
	if err := h.parseAndValidateRequest(c, &request); err != nil {
		return
	}

	// Execute the usecase
	tenant, err := h.registerUsecase.Execute(c, request)
	if err != nil {
		c.Error(custom_errors.New(err, custom_errors.TenantUnprocessEntity, "failed to register the tenant."))
		return
	}

	// Transform the tenant entity into a response DTO
	response := NewRegisterTenantResponse(tenant)

	// Respond with the created tenant
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Tenant registered successfully",
		"data":    response,
	})
}

func (h *tenantRegisterHandler) parseAndValidateRequest(c *gin.Context, request *RegisterTenantRequest) error {
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
		validationErrors := request.CustomValidationMessages(err)
		c.Error(custom_errors.New(
			err,
			custom_errors.TenantUnprocessEntity,
			"failed to validate request payload",
			validationErrors))
		return err
	}

	return nil
}
