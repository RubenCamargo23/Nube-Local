package tasks

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) FindByVinculacionSemana(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) ResumenByUsuario(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) EstadoSemanal(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) Create(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) FindByVinculacion(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) FindByID(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) Update(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) Delete(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) AddAttachment(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) FindByMe(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
func (h *Handler) HistorialByMe(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"error": "no implementado"}) }
