package spaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"})
}

func (h *Handler) FindByProfesor(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"})
}

func (h *Handler) FindByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"})
}

func (h *Handler) Close(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"})
}
