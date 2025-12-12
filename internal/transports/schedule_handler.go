package transports

import (
	"net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type ScheduleHandler struct {
	schedule services.ScheduleService
}

func NewScheduleHandler(schedule services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		schedule: schedule,
	}
}

func (h *ScheduleHandler) RegisterRoutes(r *gin.RouterGroup) {
	s := r.Group("/schedules")
	{
		admin := s.Group("")
		admin.Use(RequireRole("admin"))
		admin.POST("", h.CreateSchedule)
		admin.GET("", h.GetSchedules)
		// admin.GET("/:id", h.GetScheduleByDoctorID)
		// admin.GET("/:id", h.GetScheduleByID)
		// admin.PATCH("/:id", h.UpdateSchedule)
		// admin.DELETE("/:id", h.DeleteSchedule)
	}

	// r.GET("/doctors/:id/schedules/available", h.GetAvailableSlots)
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req models.ScheduleCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.schedule.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

func (h *ScheduleHandler) GetSchedules(c *gin.Context) {

	schedules, err := h.schedule.ListSchedules(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedules)
}
