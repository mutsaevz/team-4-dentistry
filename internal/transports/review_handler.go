package transports

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type ReviewHandler struct {
	review services.ReviewService
	logger *slog.Logger
}

func NewReviewHandler(reviewService services.ReviewService, logger *slog.Logger) *ReviewHandler {
	return &ReviewHandler{
		review: reviewService,
		logger: logger,
	}
}

func (h *ReviewHandler) RegisterRoutes(r *gin.RouterGroup) {
	review := r.Group("/reviews")
	{
		//------user---------
		review.POST("", h.CreateReview)
		review.PUT("/:id", h.UpdateReview)
		review.DELETE("/:id", h.DeleteReview)

		review.GET("/doctor/:id", h.GetDoctorReviews)

		//-------admin---------
		admin := review.Group("")
		admin.Use(RequireRole("admin"))

		admin.GET("/:id", h.GetReviewByID)
		admin.GET("/patient/:patient_id", h.GetPatientReviews)
	}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var input models.ReviewCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	review, err := h.review.CreateReview(c.Request.Context(), input)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *ReviewHandler) GetReviewByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.review.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid review ID"})
		return
	}

	var input models.ReviewUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	review, err := h.review.UpdateReview(c.Request.Context(), uint(id), input)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.review.DeleteReview(c.Request.Context(), uint(id)); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete review"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ReviewHandler) GetDoctorReviews(c *gin.Context) {
	param := c.Param("id")
	doctorID, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid doctor ID"})
		return
	}

	reviews, err := h.review.GetDoctorReviews(c.Request.Context(), uint(doctorID))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) GetPatientReviews(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviews, err := h.review.GetPatientReviews(c.Request.Context(), uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
