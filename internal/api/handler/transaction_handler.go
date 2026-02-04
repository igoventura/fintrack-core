package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-api/internal/api/dto"
	"github.com/igoventura/fintrack-api/internal/service"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Create handles the creation of a new transaction.
// @Summary Create a new transaction
// @Description Creates a new transaction for the authenticated user's tenant.
// @Tags transactions
// @Accept json
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param transaction body dto.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transactions [post]
func (h *TransactionHandler) Create(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := req.ToDomain()
	if err := h.service.Create(c.Request.Context(), tx, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.FromTransactionDomain(tx))
}

// GetByID returns a transaction by ID.
// @Summary Get transaction by ID
// @Description Retrieves a transaction by its ID.
// @Tags transactions
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Transaction ID"
// @Success 200 {object} dto.TransactionResponse
// @Failure 404 {object} map[string]string "Transaction not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	tx, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tx == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, dto.FromTransactionDomain(tx))
}

// List returns a list of transactions with optional filtering.
// @Summary List transactions
// @Description Lists transactions for the tenant, optionally filtered.
// @Tags transactions
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param accrual_month query string false "Accrual Month (YYYYMM)"
// @Param account_id query string false "Account ID"
// @Param transaction_type query string false "Transaction Type"
// @Success 200 {array} dto.TransactionResponse
// @Failure 400 {object} map[string]string "Invalid query params"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transactions [get]
func (h *TransactionHandler) List(c *gin.Context) {
	var filterReq dto.TransactionFilterRequest
	if err := c.ShouldBindQuery(&filterReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txs, err := h.service.List(c.Request.Context(), filterReq.ToDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.TransactionResponse
	for _, tx := range txs {
		response = append(response, dto.FromTransactionDomain(&tx))
	}
	// Return empty array instead of null if no results
	if response == nil {
		response = []dto.TransactionResponse{}
	}

	c.JSON(http.StatusOK, response)
}

// Update updates an existing transaction.
// @Summary Update a transaction
// @Description Updates a transaction by ID.
// @Tags transactions
// @Accept json
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Transaction ID"
// @Param transaction body dto.CreateTransactionRequest true "Transaction data"
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transactions/{id} [put]
func (h *TransactionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := req.ToDomain()
	tx.ID = id

	// Note: Update logic in Service might need refinement for tags update (currently pending in roadmap),
	// but we pass TagIDs anyway if service handles it or ignores it.
	// Current Service.Update implementation does NOT handle tags yet (Step 703: "Update tags? ... I'll leave Update basic for now").
	// So we just call Update.
	if err := h.service.Update(c.Request.Context(), tx, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.FromTransactionDomain(tx))
}

// Delete removes a transaction.
// @Summary Delete a transaction
// @Description Soft-deletes a transaction by ID.
// @Tags transactions
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Transaction ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
