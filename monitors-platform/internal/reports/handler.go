package reports

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    jwtpkg "github.com/org/monitors-platform/pkg/jwt"
    "github.com/org/monitors-platform/pkg/response"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Generate(c *gin.Context) {
    var req GenerateReportRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, http.StatusBadRequest, err.Error())
        return
    }

    res, err := h.svc.Generate(c.Request.Context(), req)
    if err != nil {
        response.Fail(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.OK(c, http.StatusCreated, res)
}

func (h *Handler) FindByProfesor(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        response.Fail(c, http.StatusUnauthorized, "no auth")
        return
    }
    userClaims := claims.(*jwtpkg.Claims)

    reports, err := h.svc.FindByProfesor(c.Request.Context(), userClaims.UserID)
    if err != nil {
        response.Fail(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, http.StatusOK, reports)
}

func (h *Handler) Download(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        response.Fail(c, http.StatusBadRequest, "id inválido")
        return
    }

    report, err := h.svc.FindByID(c.Request.Context(), id)
    if err != nil {
        response.Fail(c, http.StatusNotFound, "reporte no encontrado")
        return
    }

    if report.EstadoGeneracion != "DONE" || report.RutaPDF == "" {
        response.Fail(c, http.StatusConflict, "el pdf no está listo")
        return
    }

    c.File(report.RutaPDF)
}
