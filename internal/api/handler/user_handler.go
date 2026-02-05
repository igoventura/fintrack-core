package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-api/domain"
	"github.com/igoventura/fintrack-api/internal/api/dto"
	"github.com/igoventura/fintrack-api/internal/service"
)

type UserHandler struct {
	userService         *service.UserService
	supabaseAuthService *service.SupabaseAuthService
}

func NewUserHandler(userService *service.UserService, supabaseAuthService *service.SupabaseAuthService) *UserHandler {
	return &UserHandler{userService: userService, supabaseAuthService: supabaseAuthService}
}

// GetProfile returns the profile of the authenticated user
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags users
// @Produce json
// @Security AuthPassword
// @Success 200 {object} dto.UserResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := domain.GetUserID(c.Request.Context())
	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// UpdateProfile updates the profile of the authenticated user
// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security AuthPassword
// @Param request body dto.UpdateUserRequest true "Update User Request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// get req
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := domain.GetUserID(c.Request.Context())

	// get original user
	originalUser, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// map req to user
	user := &domain.User{
		ID:         userID,
		Name:       req.Name,
		Email:      req.Email,
		SupabaseID: originalUser.SupabaseID,
	}

	// validate user
	isValid, validationError := user.IsValid()
	if !isValid {
		jsonString, _ := json.Marshal(validationError)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: string(jsonString)})
		return
	}

	// update user
	if err := h.userService.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// update supabase user
	if originalUser.Email != req.Email || originalUser.Name != req.Name {
		if err := h.supabaseAuthService.UpdateUser(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// ListUserTenants returns the tenants of the authenticated user
// @Summary List user tenants
// @Description Get the tenants of the authenticated user
// @Tags users
// @Produce json
// @Security AuthPassword
// @Success 200 {object} []dto.UserTenantResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/tenants [get]
func (h *UserHandler) ListUserTenants(c *gin.Context) {
	userID := domain.GetUserID(c.Request.Context())
	tenants, err := h.userService.ListUserTenants(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	var resp []*dto.UserTenantResponse
	for _, tenant := range tenants {
		resp = append(resp, &dto.UserTenantResponse{
			TenantID: tenant.ID,
			Name:     tenant.Name,
		})
	}

	c.JSON(http.StatusOK, resp)
}
