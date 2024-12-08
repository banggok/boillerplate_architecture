package register

import (
	"appointment_management_system/internal/domain/tenants/rest/v1/register/dto"
	"appointment_management_system/internal/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TenantRegisterHandler represents the HTTP handler for tenant-related operations
type TenantRegisterHandler interface {
	Register(*gin.Context)
}

type tenantRegisterHandler struct {
	validator       *validator.Validate
	registerUsecase TenantRegisterUsecase
}

// NewTenantRegisterHandler creates a new TenantHandler instance with a validator
func newTenantRegisterHandler(
	v *validator.Validate,
	uc TenantRegisterUsecase,
) TenantRegisterHandler {
	return &tenantRegisterHandler{
		validator:       v,
		registerUsecase: uc,
	}
}

// Register handles tenant registration
func (h *tenantRegisterHandler) Register(c *gin.Context) {
	var request dto.RegisterTenantRequest

	// Parse and validate the request payload
	if err := h.parseAndValidateRequest(c, &request); err != nil {
		return
	}

	// Execute the usecase
	tenant, err := h.registerUsecase.Execute(c, request)
	if err != nil {
		h.handleError(
			c,
			http.StatusUnprocessableEntity,
			"TENANT_REGISTRATION_FAILED",
			"Failed to register the tenant.",
			map[string]string{
				"error": err.Error(),
			},
		)
		return
	}

	// Transform the tenant entity into a response DTO
	response := dto.NewRegisterTenantResponse(tenant)

	// Respond with the created tenant
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Tenant registered successfully",
		"data":    response,
	})
}

func (h *tenantRegisterHandler) parseAndValidateRequest(c *gin.Context, request *dto.RegisterTenantRequest) error {
	// Parse the incoming JSON request
	if err := c.ShouldBindJSON(request); err != nil {
		h.handleError(
			c,
			http.StatusBadRequest,
			"INVALID_PAYLOAD",
			"Invalid request payload",
			map[string]string{
				"error": err.Error(),
			})
		return err
	}

	// Validate the request using the validator
	if err := h.validator.Struct(request); err != nil {
		validationErrors := request.CustomValidationMessages(err)
		h.handleError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", validationErrors)
		return err
	}

	return nil
}

func (h *tenantRegisterHandler) handleError(c *gin.Context, status int, code, message string, details map[string]string) {
	c.Error(&middleware.CustomError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	})
}
