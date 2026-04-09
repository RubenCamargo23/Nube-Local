package assignments

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) FindByEspacio(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) FindByProfesor(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) Update(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) FindByMe(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
