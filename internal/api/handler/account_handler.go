package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-api/domain"
	"github.com/igoventura/fintrack-api/internal/api/dto"
	"github.com/igoventura/fintrack-api/internal/service"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(s *service.AccountService) *AccountHandler {
	return &AccountHandler{service: s}
}

// Get godoc
// @Summary Get an account
// @Description get account by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Security AuthPassword
// @Success 200 {object} dto.AccountResponse
// @Failure 404 {object} handler.ErrorResponse
// @Router /accounts/{id} [get]
func (h *AccountHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ErrorJSON(c, http.StatusBadRequest, "Account ID is required")
		return
	}

	acc, err := h.service.GetAccount(c.Request.Context(), id)
	if err != nil {
		ErrorJSON(c, http.StatusNotFound, "Account not found")
		return
	}

	c.JSON(http.StatusOK, dto.MapAccountToResponse(acc))
}

// List godoc
// @Summary List accounts
// @Description list accounts by tenant ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Security AuthPassword
// @Success 200 {array} dto.AccountResponse
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounts [get]
func (h *AccountHandler) List(c *gin.Context) {
	accounts, err := h.service.ListAccounts(c.Request.Context())
	if err != nil {
		ErrorJSON(c, http.StatusInternalServerError, "Failed to list accounts")
		return
	}

	var resp []*dto.AccountResponse
	for _, a := range accounts {
		resp = append(resp, dto.MapAccountToResponse(&a))
	}

	c.JSON(http.StatusOK, resp)
}

// Create godoc
// @Summary Create an account
// @Description create a new account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param account body dto.CreateAccountRequest true "Create account"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Security AuthPassword
// @Success 201 {object} dto.AccountResponse
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounts [post]
func (h *AccountHandler) Create(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorJSON(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tenantID := domain.GetTenantID(c.Request.Context())
	userId := domain.GetUserID(c.Request.Context())
	if userId == "" {
		ErrorJSON(c, http.StatusBadRequest, "User not found")
		return
	}

	acc := req.ToEntity(userId, tenantID)

	if err := h.service.CreateAccount(c.Request.Context(), acc); err != nil {
		ErrorJSON(c, http.StatusInternalServerError, "Failed to create account")
		return
	}

	c.JSON(http.StatusCreated, dto.MapAccountToResponse(acc))
}

// Update godoc
// @Summary Update an account
// @Description update an existing account (cannot change type or currency)
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Param account body dto.UpdateAccountRequest true "Update account"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Security AuthPassword
// @Success 200 {object} dto.AccountResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Router /accounts/{id} [put]
func (h *AccountHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ErrorJSON(c, http.StatusBadRequest, "Account ID is required")
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorJSON(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userId := domain.GetUserID(c.Request.Context())
	if userId == "" {
		ErrorJSON(c, http.StatusBadRequest, "User not found")
		return
	}

	acc := req.ToEntity(id, userId)

	if err := h.service.UpdateAccount(c.Request.Context(), acc); err != nil {
		ErrorJSON(c, http.StatusInternalServerError, "Failed to update account")
		return
	}

	c.JSON(http.StatusOK, dto.MapAccountToResponse(acc))
}

// Delete godoc
// @Summary Delete an account
// @Description delete an account by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Security AuthPassword
// @Success 204 "No Content"
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Router /accounts/{id} [delete]
func (h *AccountHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ErrorJSON(c, http.StatusBadRequest, "Account ID is required")
		return
	}

	userId := domain.GetUserID(c.Request.Context())

	if err := h.service.DeleteAccount(c.Request.Context(), id, userId); err != nil {
		ErrorJSON(c, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	c.Status(http.StatusNoContent)
}
