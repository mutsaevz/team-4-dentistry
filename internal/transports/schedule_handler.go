package transports

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type ScheduleHandler struct {
	schedule services.ScheduleService
	logger   *slog.Logger
}

func NewScheduleHandler(schedule services.ScheduleService, logger *slog.Logger) *ScheduleHandler {
	return &ScheduleHandler{
		schedule: schedule,
		logger:   logger,
	}
}

func (h *ScheduleHandler) RegisterRoutes(r *gin.RouterGroup) {
	s := r.Group("/schedules")
	{
		admin := s.Group("")
		admin.Use(RequireRole("admin"))
		admin.POST("", h.CreateSchedule)
		admin.GET("", h.GetSchedules)
		admin.GET("/:id", h.GetScheduleByDoctorID)
		admin.PATCH("/:id", h.UpdateSchedule)
		admin.DELETE("/:id", h.DeleteSchedule)
	}
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req models.ScheduleCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Ошибка разбора тела запроса (CreateSchedule)", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Debug("Запрос на создание расписания", "doctor_id", req.DoctorID, "date", req.Date)

	schedule, err := h.schedule.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Не удалось создать расписание", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Создано расписание", "schedule_id", schedule.ID, "doctor_id", schedule.DoctorID)
	c.JSON(http.StatusCreated, schedule)
}

func (h *ScheduleHandler) GetSchedules(c *gin.Context) {

	h.logger.Debug("Запрос списка расписаний")

	schedules, err := h.schedule.ListSchedules(c.Request.Context())

	if err != nil {
		h.logger.Error("Не удалось получить список расписаний", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Список расписаний получен", "count", len(schedules))
	c.JSON(http.StatusOK, schedules)
}

func (h *ScheduleHandler) GetScheduleByDoctorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Ошибка парсинга ID врача (GetScheduleByDoctorID)", "error", err.Error(), "param", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Debug("Запрос расписания по doctor_id", "doctor_id", id)

	schedules, err := h.schedule.GetSchedulesByID(c.Request.Context(), uint(id))

	if err != nil {
		h.logger.Error("Не удалось получить расписание по doctor_id", "error", err.Error(), "doctor_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Расписание получено по doctor_id", "doctor_id", id, "count", len(schedules))
	c.JSON(http.StatusOK, schedules)
}

func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Warn("Ошибка парсинга ID расписания (UpdateSchedule)", "error", err.Error(), "param", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.ScheduleUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Ошибка разбора тела запроса (UpdateSchedule)", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Debug("Запрос на обновление расписания", "schedule_id", id)

	update, err := h.schedule.UpdateSchedule(c.Request.Context(), uint(id), req)
	if err != nil {
		h.logger.Error("Не удалось обновить расписание", "error", err.Error(), "schedule_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Расписание успешно обновлено", "schedule_id", update.ID)
	c.JSON(http.StatusOK, update)
}

func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Warn("Ошибка парсинга ID расписания (DeleteSchedule)", "error", err.Error(), "param", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Debug("Запрос на удаление расписания", "schedule_id", id)

	if err := h.schedule.DeleteSchedule(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Не удалось удалить расписание", "error", err.Error(), "schedule_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Расписание удалено", "schedule_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "schedule deleted"})
}
