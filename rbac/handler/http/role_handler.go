package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"github.com/rizkiramadhan2/rbac/internal/entity"
	"github.com/rizkiramadhan2/rbac/internal/usecase"
)

type RoleHandler struct {
	roleUsecase *usecase.RoleUsecase
}

func NewRoleHandler(roleUsecase *usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{roleUsecase: roleUsecase}
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req entity.Role

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.roleUsecase.Create(c.Request.Context(), req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusCreated, gin.H{"message": "role created"})
}

func (h *RoleHandler) Get(c *gin.Context) {
	role, err := h.roleUsecase.Get(
		c.Request.Context(),
		c.Param("workspace_id"),
		c.Param("role_id"),
	)
	if err != nil {
		c.JSON(nethttp.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, role)
}

func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.roleUsecase.List(c.Request.Context(), c.Param("workspace_id"))
	if err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, roles)
}

func (h *RoleHandler) AssignToUser(c *gin.Context) {
	var req entity.UserRoleAssignment

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.roleUsecase.AssignToUser(c.Request.Context(), req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"message": "role assigned to user"})
}

func (h *RoleHandler) RemoveFromUser(c *gin.Context) {
	err := h.roleUsecase.RemoveFromUser(
		c.Request.Context(),
		c.Param("workspace_id"),
		c.Param("user_id"),
		c.Param("role_id"),
	)
	if err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{"message": "role removed from user"})
}
