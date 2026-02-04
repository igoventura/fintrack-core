package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/domain"
	"github.com/igoventura/fintrack-core/internal/api/dto"
	"github.com/igoventura/fintrack-core/internal/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategory creates a new category
// @Summary Create category
// @Description Create a new category for the authenticated user's tenant
// @Tags categories
// @Accept json
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body dto.CreateCategoryRequest true "Create Category Request"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := domain.GetUserID(c.Request.Context())
	category := &domain.Category{
		Name:             req.Name,
		ParentCategoryID: req.ParentCategoryID,
		Color:            req.Color,
		Icon:             req.Icon,
		CreatedBy:        userID,
		UpdatedBy:        userID,
	}

	if err := h.service.CreateCategory(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CategoryResponse{
		ID:               category.ID,
		ParentCategoryID: category.ParentCategoryID,
		TenantID:         category.TenantID,
		Name:             category.Name,
		Color:            category.Color,
		Icon:             category.Icon,
		CreatedAt:        category.CreatedAt,
		CreatedBy:        category.CreatedBy,
		UpdatedAt:        category.UpdatedAt,
		UpdatedBy:        category.UpdatedBy,
	})
}

// GetCategory returns a category by ID
// @Summary Get category
// @Description Get a category by its ID
// @Tags categories
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Category ID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.service.GetCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CategoryResponse{
		ID:               category.ID,
		ParentCategoryID: category.ParentCategoryID,
		TenantID:         category.TenantID,
		Name:             category.Name,
		Color:            category.Color,
		Icon:             category.Icon,
		CreatedAt:        category.CreatedAt,
		CreatedBy:        category.CreatedBy,
		UpdatedAt:        category.UpdatedAt,
		UpdatedBy:        category.UpdatedBy,
	})
}

// ListCategories returns all categories for the tenant
// @Summary List categories
// @Description Get all categories for the authenticated user's tenant
// @Tags categories
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} []dto.CategoryResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	var resp []*dto.CategoryResponse
	for _, category := range categories {
		resp = append(resp, &dto.CategoryResponse{
			ID:               category.ID,
			ParentCategoryID: category.ParentCategoryID,
			TenantID:         category.TenantID,
			Name:             category.Name,
			Color:            category.Color,
			Icon:             category.Icon,
			CreatedAt:        category.CreatedAt,
			CreatedBy:        category.CreatedBy,
			UpdatedAt:        category.UpdatedAt,
			UpdatedBy:        category.UpdatedBy,
		})
	}

	if resp == nil {
		resp = []*dto.CategoryResponse{}
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateCategory updates a category
// @Summary Update category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Update Category Request"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := domain.GetUserID(c.Request.Context())
	category := &domain.Category{
		ID:               id,
		Name:             req.Name,
		ParentCategoryID: req.ParentCategoryID,
		Color:            req.Color,
		Icon:             req.Icon,
		UpdatedBy:        userID,
	}

	if err := h.service.UpdateCategory(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	updatedCategory, err := h.service.GetCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CategoryResponse{
		ID:               updatedCategory.ID,
		ParentCategoryID: updatedCategory.ParentCategoryID,
		TenantID:         updatedCategory.TenantID,
		Name:             updatedCategory.Name,
		Color:            updatedCategory.Color,
		Icon:             updatedCategory.Icon,
		CreatedAt:        updatedCategory.CreatedAt,
		CreatedBy:        updatedCategory.CreatedBy,
		UpdatedAt:        updatedCategory.UpdatedAt,
		UpdatedBy:        updatedCategory.UpdatedBy,
	})
}

// DeleteCategory deletes a category
// @Summary Delete category
// @Description Soft Delete a category
// @Tags categories
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Category ID"
// @Success 204 "No Content"
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	userID := domain.GetUserID(c.Request.Context())

	if err := h.service.DeleteCategory(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
