package assignments

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

func (h *Handler) Create(c *gin.Context) {
	espacioIDStr := c.Param("id")
	espacioID, err := uuid.Parse(espacioIDStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id de espacio inválido")
		return
	}

	claims, exists := c.Get("claims")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "no auth")
		return
	}
	userClaims := claims.(*jwtpkg.Claims)

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.Create(c.Request.Context(), espacioID, userClaims.UserID, req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusCreated, res)
}

func (h *Handler) FindByEspacio(c *gin.Context) {
	espacioIDStr := c.Param("id")
	espacioID, err := uuid.Parse(espacioIDStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id de espacio inválido")
		return
	}

	res, err := h.svc.FindByEspacio(c.Request.Context(), espacioID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusOK, res)
}

func (h *Handler) FindByProfesor(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "no auth")
		return
	}
	userClaims := claims.(*jwtpkg.Claims)

	res, err := h.svc.FindByProfesor(c.Request.Context(), userClaims.UserID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
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

	var req struct {
		HorasSemanales int `json:"horas_semanales" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Update(c.Request.Context(), id, req.HorasSemanales); err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusOK, gin.H{"mensaje": "vinculación actualizada"})
}

func (h *Handler) FindByMe(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "no auth")
		return
	}
	userClaims := claims.(*jwtpkg.Claims)

	res, err := h.svc.FindByMe(c.Request.Context(), userClaims.UserID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusOK, res)
}
