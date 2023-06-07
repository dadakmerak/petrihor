package routes

import (
	"net/http"

	"github.com/dadakmerak/petrihor/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	v1.GET("/", h.list)
	v1.GET("/:id", h.detail)
	return router
}

func (h *Handler) list(c *gin.Context) {
	lists, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, lists)
}

func (h *Handler) detail(c *gin.Context) {
	detail, err := h.service.Detail()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, detail)
}
