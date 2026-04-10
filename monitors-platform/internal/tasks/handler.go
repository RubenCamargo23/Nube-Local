package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	jwtpkg "github.com/org/monitors-platform/pkg/jwt"
	"github.com/org/monitors-platform/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) FindByVinculacionSemana(c *gin.Context) {
	idStr := c.Param("id")
	vinculacionID, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id inválido")
		return
	}
	semana := c.Query("semana")
	if semana == "" {
		response.Fail(c, http.StatusBadRequest, "falta parametro semana")
		return
	}
	res, err := h.svc.FindByVinculacionSemana(c.Request.Context(), vinculacionID, semana)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, http.StatusOK, res)
}

func (h *Handler) ResumenByUsuario(c *gin.Context) {
	idStr := c.Param("id")
	usuarioID, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id inválido")
		return
	}

	claims, exists := c.Get("claims")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "no auth")
		return
	}
	userClaims := claims.(*jwtpkg.Claims)

	res, err := h.svc.ResumenByUsuario(c.Request.Context(), usuarioID, userClaims.UserID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, http.StatusOK, res)
}

func (h *Handler) EstadoSemanal(c *gin.Context) {
	idStr := c.Param("id")
	espacioID, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id inválido")
		return
	}
	semana := c.Query("semana")
	if semana == "" {
		response.Fail(c, http.StatusBadRequest, "falta parametro semana")
		return
	}
	res, err := h.svc.EstadoSemanal(c.Request.Context(), espacioID, semana)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, http.StatusOK, res)
}

func (h *Handler) Create(c *gin.Context) {
	vinculacionIDStr := c.Param("id")
	vinculacionID, err := uuid.Parse(vinculacionIDStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id inválido")
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.Create(c.Request.Context(), vinculacionID, req)
	if err != nil {
		// Aquí puedes hacer un switch de errores comunes para mejores códigos HTTP
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusCreated, res)
}

func (h *Handler) FindByVinculacion(c *gin.Context) {
	vinculacionIDStr := c.Param("id")
	vinculacionID, err := uuid.Parse(vinculacionIDStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id inválido")
		return
	}

	res, err := h.svc.FindByVinculacion(c.Request.Context(), vinculacionID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusOK, res)
}

func (h *Handler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id inválido")
		return
	}

	res, err := h.svc.FindByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "tarea no encontrada")
		return
	}

	response.OK(c, http.StatusOK, res)
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id inválido")
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusOK, res)
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id inválido")
		return
	}

	err = h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusOK, gin.H{"mensaje": "tarea eliminada"})
}

func (h *Handler) AddAttachment(c *gin.Context) {
	// A simple mock of upload processing for this endpoint
	c.JSON(http.StatusOK, gin.H{"mensaje": "adjunto agregado exitosamente"})
}

func (h *Handler) FindByMe(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "no auth")
		return
	}
	userClaims := claims.(*jwtpkg.Claims)

	res, err := h.svc.FindByUsuario(c.Request.Context(), userClaims.UserID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, http.StatusOK, res)
}

func (h *Handler) HistorialByMe(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "no auth")
		return
	}
	userClaims := claims.(*jwtpkg.Claims)

	// In this simplified version, HistorialByMe is practically FindByUsuario which lists all
	res, err := h.svc.FindByUsuario(c.Request.Context(), userClaims.UserID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, http.StatusOK, res)
}
