package transports

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type DoctorHandler struct {
	doctor   services.DoctorService
	service  services.ServService
	schedule services.ScheduleService
	review   services.ReviewService
	logger *slog.Logger
}

func NewDoctorHandler(
	doctor services.DoctorService,
	service services.ServService,
	schedule services.ScheduleService,
	review services.ReviewService,
	logger *slog.Logger) *DoctorHandler {
	return &DoctorHandler{
		doctor:   doctor,
		service:  service,
		schedule: schedule,
		review:   review,
		logger: logger,
	}
}

func (h *DoctorHandler) RegisterRoutes(r *gin.RouterGroup) {
	doctor := r.Group("/doctors")

	{
		//-----patient------
		doctor.GET("/:id", h.GetDoctorByID)
		doctor.GET("", h.ListDoctors)

		doctor.GET("/:id/reviews", h.GetDoctorReviews)

		doctor.GET("/:id/services", h.ListDoctorServices)

		doctor.GET("/:id/schedules/available", h.GetAvailableSlots)

		//-----admin------
		admin := doctor.Group("")
		admin.Use(RequireRole("admin"))

		admin.POST("", h.CreateDoctor)
		admin.PATCH("/:id", h.UpdateDoctor)
		admin.DELETE("/:id", h.DeleteDoctor)
		admin.GET("/:id/schedules", h.ListSchedules)
	}
}

func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var input models.DoctorCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor, err := h.doctor.CreateDoctor(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to create doctor")
		return
	}

	c.JSON(http.StatusCreated, doctor)
}

func (h *DoctorHandler) GetDoctorByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

	doctor, err := h.doctor.GetDoctorByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to get doctor")
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) ListDoctors(c *gin.Context) {

	doctors, err := h.doctor.ListDoctors(c.Request.Context(), GetDoctorQueryParams(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to list doctors")
		return
	}

	c.JSON(http.StatusOK, doctors)
}

func (h *DoctorHandler) UpdateDoctor(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

	var input models.DoctorUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor, err := h.doctor.UpdateDoctor(c.Request.Context(), uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to update doctor")
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) DeleteDoctor(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

	if err := h.doctor.DeleteDoctor(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to delete doctor")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DoctorHandler) GetDoctorReviews(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

	reviews, err := h.review.GetDoctorReviews(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to get doctor reviews")
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *DoctorHandler) ListSchedules(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

	schedules, err := h.schedule.GetSchedulesByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to list schedules")
		return
	}

	c.JSON(http.StatusOK, schedules)
}

func (h *DoctorHandler) ListDoctorServices(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

	services, err := h.doctor.GetDoctorServices(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to list doctor services")
		return
	}

	c.JSON(http.StatusOK, services)
}

func (h *DoctorHandler) GetAvailableSlots(c *gin.Context) {
	doctorID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	available, err := h.schedule.GetAvailableSlots(c.Request.Context(), uint(doctorID), QueryWeek(c))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, available)
}

func QueryWeek(c *gin.Context) int {
	week := c.Query("week")

	w, _ := strconv.Atoi(week)

	return w
}
func GetDoctorQueryParams(c *gin.Context) models.DoctorQueryParams {
	specialization := c.Query("specialization")
	experience := c.Query("experience")
	avg := c.Query("avg")

	filOr := c.Query("fil_or")

	v, _ := strconv.ParseBool(filOr)

	eID, _ := strconv.Atoi(experience)

	aID, _ := strconv.Atoi(avg)

	params := models.DoctorQueryParams{
		Specialization:  specialization,
		ExperienceYears: eID,
		AvgRating:       float64(aID),

		FilOr: v,
	}

	return params
}
