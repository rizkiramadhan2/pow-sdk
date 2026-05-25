package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkiramadhan2/rbac/internal/entity"
	"github.com/rizkiramadhan2/rbac/internal/usecase"
)

type WorkspaceHandler struct {
	workspaceUsecase *usecase.WorkspaceUsecase
}

func NewWorkspaceHandler(workspaceUsecase *usecase.WorkspaceUsecase) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceUsecase: workspaceUsecase}
}

func (h *WorkspaceHandler) Create(c *gin.Context) {
	var req entity.Workspace

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.workspaceUsecase.Create(c.Request.Context(), req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusCreated, gin.H{"message": "workspace created"})
}

func (h *WorkspaceHandler) Get(c *gin.Context) {
	workspace, err := h.workspaceUsecase.Get(c.Request.Context(), c.Param("workspace_id"))
	if err != nil {
		c.JSON(nethttp.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, workspace)
}
