package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"github.com/rizkiramadhan2/rbac/internal/entity"
	"github.com/rizkiramadhan2/rbac/internal/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req entity.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUsecase.Create(c.Request.Context(), req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusCreated, gin.H{"message": "user created"})
}

func (h *UserHandler) Get(c *gin.Context) {
	user, err := h.userUsecase.Get(
		c.Request.Context(),
		c.Param("workspace_id"),
		c.Param("user_id"),
	)
	if err != nil {
		c.JSON(nethttp.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(nethttp.StatusOK, user)
}
