package transports

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type RecommendationHandler struct {
	service services.RecommendationService
}

func NewRecommendationHandler(
	service services.RecommendationService,
) *RecommendationHandler {
	return &RecommendationHandler{service: service}
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка в данных токена",
		})
		return
	}

	if role != "doctor" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "нет прав для создания рекомендаций",
		})
		return
	}

	var req models.RecommendationCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	rec, err := h.service.CreateRec(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный userID в контексте",
		})
		return
	}

	recs, err := h.service.ListRecsByPatientID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "нет прав на удаление рекомендаций",
		})
		return
	}

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	if err := h.service.DeleteRec(uint(id)); err != nil {
		if errors.Is(err, services.ErrRecommendationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
