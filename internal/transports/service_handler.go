package transports

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type ServiceHandler struct {
	service services.ServService
	logger  *slog.Logger
}

func NewServiceHandler(
	service services.ServService,
	logger *slog.Logger,
) *ServiceHandler {
	return &ServiceHandler{service: service, logger: logger}
}

func (h *ServiceHandler) RegisterRoutes(r *gin.RouterGroup) {
	services := r.Group("/services")

	services.GET("/:id", h.GetByID)
	services.GET("", h.List)

	admin := services.Group("")
	admin.Use(RequireRole("admin"))
	admin.POST("", h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)

}

func (h *ServiceHandler) Create(c *gin.Context) {
	var req models.ServiceCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Ошибка парсинга JSON в Service.Create", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	service, err := h.service.CreateService(req)

	if err != nil {
		h.logger.Error("Ошибка создания услуги", "error", err.Error(), "name", req.Name)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Услуга создана", "service_id", service.ID, "name", service.Name)
	c.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		h.logger.Warn("Неверный id в Service.GetByID", "param", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	service, err := h.service.GetServiceByID(uint(id))

	if err != nil {
		if errors.Is(err, services.ErrServiceNotfound) {
			h.logger.Warn("Услуга не найдена", "service_id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Ошибка получения услуги", "error", err.Error(), "service_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Услуга получена", "service_id", service.ID)
	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) List(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	category := c.Query("category")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный offset"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		h.logger.Warn("Неверный limit в Service.List", "limit", limitStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный limit"})
		return
	}

	if category != "" {
		services, err := h.service.ListServicesByCategory(category, offset, limit)
		if err != nil {
			h.logger.Error("Ошибка получения услуг по категории", "error", err.Error(), "category", category)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Услуги по категории получены", "category", category, "count", len(services))
		c.JSON(http.StatusOK, services)
		return
	}

	services, err := h.service.ListServices(offset, limit)
	if err != nil {
		h.logger.Error("Ошибка получения списка услуг", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Список услуг получен", "count", len(services))
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
			h.logger.Warn("Услуга не найдена при обновлении", "service_id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Ошибка обновления услуги", "error", err.Error(), "service_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Услуга обновлена", "service_id", service.ID)
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
			h.logger.Warn("Услуга не найдена при удалении", "service_id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Ошибка удаления услуги", "error", err.Error(), "service_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Услуга удалена", "service_id", id)
	c.Status(http.StatusOK)
}
