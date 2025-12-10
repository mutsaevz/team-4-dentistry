package transports

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type ServiceHandler struct {
	service services.ServService
}

func NewServiceHandler(
	service services.ServService,
) *ServiceHandler {
	return &ServiceHandler{service: service}
}

func (h *ServiceHandler) RegisterRoutes(r *gin.Engine) {
	services := r.Group("/services")
	{
		services.GET("/:id", h.GetByID)
		services.GET("", h.List)
		services.POST("", h.Create)
		services.PATCH("/:id", h.Update)
		services.DELETE("/:id", h.Delete)
	}
}

func (h *ServiceHandler) Create(c *gin.Context) {
	var req models.ServiceCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	service, err := h.service.CreateService(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	service, err := h.service.GetServiceByID(uint(id))

	if err != nil {
		if errors.Is(err, services.ErrServiceNotfound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) List(c *gin.Context) {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	category := c.Query("category")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный offset"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный limit"})
		return
	}

	if category != "" {
		services, err := h.service.ListServicesByCategory(category, offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, services)
	}

	services, err := h.service.ListServices(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

func (h *ServiceHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	var req models.ServiceUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	service, err := h.service.UpdateService(uint(id), req)

	if err != nil {
		if errors.Is(err, services.ErrServiceNotfound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	if err := h.service.DeleteService(uint(id)); err != nil {
		if errors.Is(err, services.ErrServiceNotfound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
