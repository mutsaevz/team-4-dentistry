package transports

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type PatientRecordHandler struct {
	service services.PatientRecordService
}

func NewPatientRecordHandler(
	service services.PatientRecordService,
) *PatientRecordHandler {
	return &PatientRecordHandler{service: service}
}

func (h *PatientRecordHandler) RegisterRoutes(c *gin.RouterGroup) {
	records := c.Group("/patient-records")
	records.POST("")
	records.GET("/")
	records.GET("/:id")
	records.PATCH("/:id")
	records.DELETE("/:id")
}

func (h *PatientRecordHandler) Create(c *gin.Context) {
	var req models.PatientRecordCreate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "некорректный JSON"})
		return
	}

	record, err := h.service.Create(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, record)
}

func (h *PatientRecordHandler) GetAll(c *gin.Context) {
	records, err := h.service.GetAll()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, records)
}

func (h *PatientRecordHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "некорректный ID"})
		return
	}

	record, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, record)
}

func (h *PatientRecordHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "некорректный ID"})
		return
	}

	var req models.PatientRecordUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "некорректный JSON"})
		return
	}

	if err := h.service.Update(uint(id), &req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "запись успешно обновлена"})
}

func (h *PatientRecordHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "некорректный ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "запись успешно удалена"})
}
