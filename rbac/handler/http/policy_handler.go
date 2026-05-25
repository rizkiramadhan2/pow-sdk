package http

import (
    nethttp "net/http"

    "github.com/gin-gonic/gin"

    "github.com/rizkiramadhan2/rbac/internal/entity"
    "github.com/rizkiramadhan2/rbac/internal/usecase"
)

type PolicyHandler struct {
    policyUsecase *usecase.PolicyUsecase
}

func NewPolicyHandler(policyUsecase *usecase.PolicyUsecase) *PolicyHandler {
    return &PolicyHandler{policyUsecase: policyUsecase}
}

func (h *PolicyHandler) Create(c *gin.Context) {
    var req entity.Policy

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.policyUsecase.Create(c.Request.Context(), req); err != nil {
        c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(nethttp.StatusCreated, gin.H{"message": "policy created"})
}

func (h *PolicyHandler) Get(c *gin.Context) {
    policy, err := h.policyUsecase.Get(
        c.Request.Context(),
        c.Param("workspace_id"),
        c.Param("policy_id"),
    )
    if err != nil {
        c.JSON(nethttp.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(nethttp.StatusOK, policy)
}

func (h *PolicyHandler) List(c *gin.Context) {
    policies, err := h.policyUsecase.List(c.Request.Context(), c.Param("workspace_id"))
    if err != nil {
        c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(nethttp.StatusOK, policies)
}

func (h *PolicyHandler) AttachToRole(c *gin.Context) {
    var req entity.RolePolicy

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.policyUsecase.AttachToRole(c.Request.Context(), req); err != nil {
        c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(nethttp.StatusOK, gin.H{"message": "policy attached to role"})
}

func (h *PolicyHandler) DetachFromRole(c *gin.Context) {
    err := h.policyUsecase.DetachFromRole(
        c.Request.Context(),
        c.Param("workspace_id"),
        c.Param("role_id"),
        c.Param("policy_id"),
    )
    if err != nil {
        c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(nethttp.StatusOK, gin.H{"message": "policy detached from role"})
}