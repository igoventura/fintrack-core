package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/domain"
	"github.com/igoventura/fintrack-core/internal/api/dto"
	"github.com/igoventura/fintrack-core/internal/service"
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

	// map req to user
	userID := domain.GetUserID(c.Request.Context())
	user := &domain.User{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
	}

	// validate user
	isValid, validationError := user.IsValid()
	if !isValid {
		jsonString, _ := json.Marshal(validationError)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: string(jsonString)})
		return
	}

	// get original user
	originalUser, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
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
