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

type RecommendationHandler struct {
	service services.RecommendationService
	logger  *slog.Logger
}

func NewRecommendationHandler(
	service services.RecommendationService,
	logger *slog.Logger,
) *RecommendationHandler {
	return &RecommendationHandler{service: service, logger: logger}
}

func (h *RecommendationHandler) RegisterRoutes(c *gin.RouterGroup) {
	recs := c.Group("/recommendations")

	recs.POST("", h.Create)

	recs.GET("/my", h.ListMy)

	recs.DELETE("/:id", h.Delete)
}

func (h *RecommendationHandler) Create(c *gin.Context) {
	idVal, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неавторизован"})
		return
	}

	roleVal, existsRole := c.Get("userRole")
	if !existsRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет роли в токене"})
		return
	}

	userID, okID := idVal.(uint)
	role, okRole := roleVal.(string)

	if !okID || !okRole {
		h.logger.Error("Некорректные данные токена у Recommendation.Create", "userID_exists", exists, "role_exists", existsRole)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка в данных токена",
		})
		return
	}

	if role != "doctor" && role != "admin" {
		h.logger.Warn("Попытка создания рекомендации без прав", "user_id", userID, "role", role)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "нет прав для создания рекомендаций",
		})
		return
	}

	var req models.RecommendationCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Ошибка парсинга JSON в Recommendation.Create", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	rec, err := h.service.CreateRec(userID, req)
	if err != nil {
		h.logger.Error("Ошибка создания рекомендации", "error", err.Error(), "doctor_id", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Рекомендация создана", "id", rec.ID, "doctor_id", userID)
	c.JSON(http.StatusCreated, rec)
}

func (h *RecommendationHandler) ListMy(c *gin.Context) {
	idVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	userID, okID := idVal.(uint)
	if !okID {
		h.logger.Error("Некорректный userID в Recommendation.ListMy", "exists", exists)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный userID в контексте",
		})
		return
	}

	recs, err := h.service.ListRecsByPatientID(userID)
	if err != nil {
		h.logger.Error("Ошибка получения рекомендаций", "error", err.Error(), "patient_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Список рекомендаций получен", "patient_id", userID, "count", len(recs))
	c.JSON(http.StatusOK, recs)
}

func (h *RecommendationHandler) Delete(c *gin.Context) {
	roleVal, exists := c.Get("userRole")

	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "нет роли в токене",
		})
		return
	}

	role, okRole := roleVal.(string)
	if !okRole || (role != "doctor" && role != "admin") {
		h.logger.Warn("Попытка удаления рекомендации без прав", "role", role)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "нет прав на удаление рекомендаций",
		})
		return
	}

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		h.logger.Warn("Неверный id в Recommendation.Delete", "param", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	if err := h.service.DeleteRec(uint(id)); err != nil {
		if errors.Is(err, services.ErrRecommendationNotFound) {
			h.logger.Warn("Recommendation не найден при удалении", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Ошибка удаления recommendation", "error", err.Error(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Recommendation удалён", "id", id)
	c.Status(http.StatusOK)
}
