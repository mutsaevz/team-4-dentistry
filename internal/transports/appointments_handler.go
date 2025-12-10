package transports

import (
	
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type AppointmentsHandler struct {
	service services.AppointmentService
}

func NewAppointmentsHandler(appointmentService services.AppointmentService) *AppointmentsHandler {
	return &AppointmentsHandler{service: appointmentService}
}

func (h *AppointmentsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	appointments := rg.Group("/appointments")
	{
		appointments.POST("/", h.Create)
		appointments.GET("/", h.GetAll)
		appointments.GET("/:id", h.GetByID)
		appointments.PATCH("/:id", h.Update)
		appointments.DELETE("/:id", h.Delete)
	}
}

func (h *AppointmentsHandler) Create(c *gin.Context) {
	var req models.AppointmentCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": constants.Invalid_JSON_Error,
		})
	}

	appointment, err := h.service.Create(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": constants.ErrCreateAppointment,
		})
		return
	}

	c.JSON(200, appointment)
}

func (h *AppointmentsHandler) GetAll(c *gin.Context) {
	appointments, err := h.service.GetAll()
	if err != nil {
		c.JSON(400, gin.H{
			"error": constants.ErrGetAppointments,
		})
	}

	c.JSON(200, appointments)
}

func (h *AppointmentsHandler) GetByID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error": constants.Parse_ID_Error,
		})
		return
	}

	appointment, err := h.service.GetByID(uint(id))
	if err != nil {

		c.JSON(404, gin.H{
			"error": constants.ErrGetByIDAppointments,
		})
		return
	}

	c.JSON(200, appointment)

}
func (h *AppointmentsHandler) Update(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
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
		c.JSON(400, gin.H{
			"error": constants.ErrUpdateAppointments,
		})
		return
	}

	c.JSON(200, gin.H{"message": "appointment updated successfully"})
}
func (h *AppointmentsHandler) Delete(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error": constants.Parse_ID_Error,
		})
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(400, gin.H{
			"error": constants.ErrDeleteAppointments,
		})
		return
	}
	
	c.JSON(200, gin.H{"message": "appointment deleted successfully"})
}
