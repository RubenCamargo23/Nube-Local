package users

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/org/monitors-platform/pkg/response"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, http.StatusBadRequest, err.Error())
        return
    }
    u, err := h.svc.Create(c.Request.Context(), req)
    if err != nil {
        response.Fail(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, http.StatusCreated, u)
}

func (h *Handler) AssignRole(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        response.Fail(c, http.StatusBadRequest, "id inválido")
        return
    }
    var req AssignRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, http.StatusBadRequest, err.Error())
        return
    }
    if err := h.svc.AssignRole(c.Request.Context(), id, req); err != nil {
        response.Fail(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, http.StatusOK, gin.H{"mensaje": "rol asignado"})
}

func (h *Handler) FindAll(c *gin.Context) {
    users, err := h.svc.FindAll(c.Request.Context())
    if err != nil {
        response.Fail(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, http.StatusOK, users)
}

func (h *Handler) FindByID(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        response.Fail(c, http.StatusBadRequest, "id inválido")
        return
    }
    u, err := h.svc.FindByID(c.Request.Context(), id)
    if err != nil {
        response.Fail(c, http.StatusNotFound, "usuario no encontrado")
        return
    }
    response.OK(c, http.StatusOK, u)
}
