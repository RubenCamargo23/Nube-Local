package periods

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/org/monitors-platform/pkg/response"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(c *gin.Context) {
    var req CreatePeriodRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, http.StatusBadRequest, err.Error())
        return
    }
    p, err := h.svc.Create(c.Request.Context(), req)
    if err != nil {
        response.Fail(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, http.StatusCreated, p)
}

func (h *Handler) FindAll(c *gin.Context) {
    periods, err := h.svc.FindAll(c.Request.Context())
    if err != nil {
        response.Fail(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, http.StatusOK, periods)
}

func (h *Handler) FindByID(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        response.Fail(c, http.StatusBadRequest, "id inválido")
        return
    }
    p, err := h.svc.FindByID(c.Request.Context(), id)
    if err != nil {
        response.Fail(c, http.StatusNotFound, "período no encontrado")
        return
    }
    response.OK(c, http.StatusOK, p)
}

func (h *Handler) Close(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        response.Fail(c, http.StatusBadRequest, "id inválido")
        return
    }
    if err := h.svc.Close(c.Request.Context(), id); err != nil {
        response.Fail(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, http.StatusOK, gin.H{"mensaje": "período cerrado"})
}
