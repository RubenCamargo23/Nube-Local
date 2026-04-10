package spaces

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
	claims, exists := c.Get("claims")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "no auth")
		return
	}
	userClaims := claims.(*jwtpkg.Claims)

	var req CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.Create(c.Request.Context(), req, userClaims.UserID)
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

	spaces, err := h.svc.FindByProfesor(c.Request.Context(), userClaims.UserID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusOK, spaces)
}

func (h *Handler) FindByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "id inválido")
		return
	}

	s, err := h.svc.FindByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "espacio no encontrado")
		return
	}

	response.OK(c, http.StatusOK, s)
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

	response.OK(c, http.StatusOK, gin.H{"mensaje": "espacio cerrado"})
}
