package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"github.com/rizkiramadhan2/rbac/internal/entity"
	"github.com/rizkiramadhan2/rbac/internal/usecase"
)

type AuthzHandler struct {
	authzUsecase *usecase.AuthzUsecase
}

func NewAuthzHandler(authzUsecase *usecase.AuthzUsecase) *AuthzHandler {
	return &AuthzHandler{authzUsecase: authzUsecase}
}

func (h *AuthzHandler) Can(c *gin.Context) {
	var req entity.CheckRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.authzUsecase.Can(c.Request.Context(), req)
	if err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, result)
}
