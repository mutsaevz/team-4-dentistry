package transports

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type AppointmentsHandler struct {
	service services.AppointmentService
	logger  *slog.Logger
}

func NewAppointmentsHandler(appointmentService services.AppointmentService, logger *slog.Logger) *AppointmentsHandler {
	return &AppointmentsHandler{service: appointmentService, logger: logger}
}

func (h *AppointmentsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	appointments := rg.Group("/appointments")

	appointments.POST("", h.Create)
	appointments.GET("/:id", h.GetByID)
	appointments.PATCH("/:id", h.Update)
	appointments.DELETE("/:id", h.Delete)
	appointments.GET("/patients/:id", h.GetByPatientID)

	admin := appointments.Group("")
	admin.Use(RequireRole("admin"))
	admin.GET("", h.GetAll)
}

func (h *AppointmentsHandler) Create(c *gin.Context) {
	var req models.AppointmentCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Ошибка парсинга JSON в Appointments.Create", "error", err.Error())
		c.JSON(400, gin.H{
			"error": constants.Invalid_JSON_Error,
		})
		return
	}
	appointment, err := h.service.Create(&req)
	if err != nil {
		h.logger.Error("Ошибка создания записи (appointment)", "error", err.Error(), "patient_id", req.PatientID)
		c.JSON(400, gin.H{
			"error": constants.ErrCreateAppointment,
		})
		return
	}

	h.logger.Info("Запись создана", "appointment_id", appointment.ID)
	c.JSON(200, appointment)
}

func (h *AppointmentsHandler) GetAll(c *gin.Context) {
	appointments, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Ошибка получения всех записей (appointments)", "error", err.Error())
		c.JSON(400, gin.H{
			"error": constants.ErrGetAppointments,
		})
		return
	}

	h.logger.Info("Список записей получен", "count", len(appointments))
	c.JSON(200, appointments)
}

func (h *AppointmentsHandler) GetByID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		h.logger.Warn("Ошибка парсинга ID в Appointments.GetByID", "param", idstr)
		c.JSON(400, gin.H{
			"error": constants.Parse_ID_Error,
		})
		return
	}

	appointment, err := h.service.GetByID(uint(id))
	if err != nil {
		h.logger.Error("Запись не найдена по ID", "error", err.Error(), "id", id)
		c.JSON(404, gin.H{
			"error": constants.ErrGetByIDAppointments,
		})
		return
	}

	h.logger.Info("Запись получена по ID", "appointment_id", appointment.ID)
	c.JSON(200, appointment)

}
func (h *AppointmentsHandler) Update(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		h.logger.Warn("Ошибка парсинга ID в Appointments.Update", "param", idstr)
		c.JSON(400, gin.H{
			"error": constants.Parse_ID_Error,
		})
		return
	}

	var req models.AppointmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": constants.Invalid_JSON_Error,
		})
		return
	}

	if err := h.service.Update(uint(id), &req); err != nil {
		h.logger.Error("Ошибка обновления записи (appointment)", "error", err.Error(), "appointment_id", id)
		c.JSON(400, gin.H{
			"error": constants.ErrUpdateAppointments,
		})
		return
	}

	h.logger.Info("Запись успешно обновлена", "appointment_id", id)
	c.JSON(200, gin.H{"message": "appointment updated successfully"})
}
func (h *AppointmentsHandler) Delete(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		h.logger.Warn("Ошибка парсинга ID в Appointments.Delete", "param", idstr)
		c.JSON(400, gin.H{
			"error": constants.Parse_ID_Error,
		})
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		h.logger.Error("Ошибка удаления записи (appointment)", "error", err.Error(), "appointment_id", id)
		c.JSON(400, gin.H{
			"error": constants.ErrDeleteAppointments,
		})
		return
	}

	h.logger.Info("Запись удалена", "appointment_id", id)
	c.JSON(200, gin.H{"message": "appointment deleted successfully"})
}

func (h *AppointmentsHandler) GetByPatientID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		h.logger.Warn("Ошибка парсинга ID в Appointments.GetByPatientID", "param", idstr)
		c.JSON(400, gin.H{
			"error": constants.Parse_ID_Error,
		})
		return
	}

	appointments, err := h.service.GetByPatientID(uint(id))
	if err != nil {
		h.logger.Error("Ошибка получения записей по patient_id", "error", err.Error(), "patient_id", id)
		c.JSON(404, gin.H{
			"error": constants.PatientIDIsIncorrect,
		})
		return
	}

	h.logger.Info("Записи пациента получены", "patient_id", id, "count", len(appointments))
	c.JSON(200, appointments)
}
