package auth

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/org/monitors-platform/pkg/response"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, http.StatusBadRequest, err.Error())
        return
    }

    res, err := h.svc.Login(c.Request.Context(), req)
    if err != nil {
        if errors.Is(err, ErrInvalidCredentials) {
            response.Fail(c, http.StatusUnauthorized, "credenciales inválidas")
            return
        }
        response.Fail(c, http.StatusInternalServerError, "error interno")
        return
    }

    response.OK(c, http.StatusOK, res)
}
