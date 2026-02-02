package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/internal/api/dto"
	"github.com/igoventura/fintrack-core/internal/service"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(s *service.AccountService) *AccountHandler {
	return &AccountHandler{service: s}
}

func (h *AccountHandler) Get(c *gin.Context) {
	id := c.Param("id")
	acc, err := h.service.GetAccount(c.Request.Context(), id)
	if err != nil {
		ErrorJSON(c, http.StatusNotFound, "Account not found")
		return
	}

	c.JSON(http.StatusOK, dto.MapAccountToResponse(acc))
}

func (h *AccountHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		ErrorJSON(c, http.StatusBadRequest, "tenant_id is required")
		return
	}

	accounts, err := h.service.ListAccounts(c.Request.Context(), tenantID)
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

func (h *AccountHandler) Create(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorJSON(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// For now, we stub the creator ID. In a real app, this would come from the auth context.
	creatorID := "00000000-0000-0000-0000-000000000001"
	acc := req.ToEntity(creatorID)

	if err := h.service.CreateAccount(c.Request.Context(), acc); err != nil {
		ErrorJSON(c, http.StatusInternalServerError, "Failed to create account")
		return
	}

	c.JSON(http.StatusCreated, dto.MapAccountToResponse(acc))
}
