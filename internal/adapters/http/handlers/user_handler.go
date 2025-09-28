package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"user-service/internal/application/dto"
	"user-service/internal/application/usecases"
	domainErrors "user-service/internal/domain/errors"
	"user-service/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUseCases usecases.UserUseCases
	validator    *validator.Validate
	logger       logger.Logger
}

func NewUserHandler(userUseCases usecases.UserUseCases, log logger.Logger) *UserHandler {
	return &UserHandler{
		userUseCases: userUseCases,
		validator:    validator.New(),
		logger:       log.With("component", "user_handler"),
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	h.logger.Info("Create user request received",
		"request_id", requestID,
		"remote_ip", c.RealIP(),
		"user_agent", c.Request().UserAgent())

	// Parse request body
	var request dto.CreateUserRequestDTO
	if err := c.Bind(&request); err != nil {
		h.logger.Warn("Failed to bind request body",
			"request_id", requestID,
			"error", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body format",
		})
	}

	// Validate request
	if err := h.validator.Struct(request); err != nil {
		h.logger.Warn("Request validation failed",
			"request_id", requestID,
			"error", err)

		details := make(map[string]interface{})
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				details[fieldError.Field()] = getValidationErrorMessage(fieldError)
			}
		}

		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "Request validation failed",
			Details: details,
		})
	}

	// Execute use case
	response, err := h.userUseCases.CreateUser(c.Request().Context(), &request)
	if err != nil {
		return h.handleError(c, err, requestID, "Failed to create user")
	}

	h.logger.Info("User created successfully",
		"request_id", requestID,
		"user_id", response.ID,
		"email", response.Email)

	return c.JSON(http.StatusCreated, response)
}

// GetUser handles GET /api/v1/users/:id
func (h *UserHandler) GetUser(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// Parse user ID from path parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid user ID parameter",
			"request_id", requestID,
			"id_param", idParam,
			"error", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid user ID format",
		})
	}

	h.logger.Info("Get user request received",
		"request_id", requestID,
		"user_id", id,
		"remote_ip", c.RealIP())

	// Execute use case
	response, err := h.userUseCases.GetUserByID(c.Request().Context(), uint(id))
	if err != nil {
		return h.handleError(c, err, requestID, "Failed to get user")
	}

	h.logger.Info("User retrieved successfully",
		"request_id", requestID,
		"user_id", response.ID)

	return c.JSON(http.StatusOK, response)
}

// GetUserByEmail handles GET /api/v1/users/email/:email
func (h *UserHandler) GetUserByEmail(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	email := c.Param("email")
	if email == "" {
		h.logger.Warn("Empty email parameter",
			"request_id", requestID)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_EMAIL",
			Message: "Email parameter is required",
		})
	}

	h.logger.Info("Get user by email request received",
		"request_id", requestID,
		"email", email,
		"remote_ip", c.RealIP())

	// Execute use case
	response, err := h.userUseCases.GetUserByEmail(c.Request().Context(), email)
	if err != nil {
		return h.handleError(c, err, requestID, "Failed to get user by email")
	}

	h.logger.Info("User retrieved by email successfully",
		"request_id", requestID,
		"user_id", response.ID,
		"email", response.Email)

	return c.JSON(http.StatusOK, response)
}

// ListUsers handles GET /api/v1/users
func (h *UserHandler) ListUsers(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	h.logger.Info("List users request received",
		"request_id", requestID,
		"remote_ip", c.RealIP())

	// Parse query parameters
	page := 1
	pageSize := 10

	if pageParam := c.QueryParam("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p >= 0 {
			page = p
		}
	}

	if sizeParam := c.QueryParam("page_size"); sizeParam != "" {
		if ps, err := strconv.Atoi(sizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	h.logger.Info("List users parameters",
		"request_id", requestID,
		"page", page,
		"page_size", pageSize)

	// Execute use case
	response, err := h.userUseCases.ListUsers(c.Request().Context(), page, pageSize)
	if err != nil {
		return h.handleError(c, err, requestID, "Failed to list users")
	}

	h.logger.Info("Users listed successfully",
		"request_id", requestID,
		"count", len(response.Users),
		"page", page)

	return c.JSON(http.StatusOK, response)
}

// handleError handles different types of errors and returns appropriate HTTP responses
func (h *UserHandler) handleError(c echo.Context, err error, requestID, logMessage string) error {
	h.logger.Error(logMessage,
		"request_id", requestID,
		"error", err)

	// Handle domain errors
	var domainErr *domainErrors.DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Code {
		case domainErrors.ErrUserNotFound.Code:
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   domainErr.Code,
				Message: domainErr.Message,
			})
		case domainErrors.ErrUserAlreadyExists.Code:
			return c.JSON(http.StatusConflict, ErrorResponse{
				Error:   domainErr.Code,
				Message: domainErr.Message,
			})
		case domainErrors.ErrInvalidUserEmail.Code,
			domainErrors.ErrInvalidUserPassword.Code:
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   domainErr.Code,
				Message: domainErr.Message,
			})
		default:
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   domainErr.Code,
				Message: domainErr.Message,
			})
		}
	}

	// Handle generic errors
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "INTERNAL_ERROR",
		Message: "An internal error occurred",
	})
}

// getValidationErrorMessage returns a user-friendly validation error message
func getValidationErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Minimum length is " + fieldError.Param() + " characters"
	case "max":
		return "Maximum length is " + fieldError.Param() + " characters"
	default:
		return "Invalid value"
	}
}
