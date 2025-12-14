package transports

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type PatientRecordHandler struct {
	service services.PatientRecordService
	logger  *slog.Logger
}

func NewPatientRecordHandler(
	service services.PatientRecordService,
	logger *slog.Logger,
) *PatientRecordHandler {
	return &PatientRecordHandler{service: service, logger: logger}
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
		h.logger.Warn("Ошибка парсинга JSON в PatientRecord.Create", "error", err.Error())
		c.JSON(400, gin.H{"error": "некорректный JSON"})
		return
	}

	record, err := h.service.Create(&req)
	if err != nil {
		h.logger.Error("Ошибка создания записи пациента", "error", err.Error(), "patient_id", req.PatientID)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Patient record создан", "id", record.ID, "patient_id", record.PatientID)
	c.JSON(200, record)
}

func (h *PatientRecordHandler) GetAll(c *gin.Context) {
	records, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Ошибка получения записей пациентов", "error", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Список patient records получен", "count", len(records))
	c.JSON(200, records)
}

func (h *PatientRecordHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Неверный ID в PatientRecord.GetByID", "param", idStr)
		c.JSON(400, gin.H{"error": "некорректный ID"})
		return
	}

	record, err := h.service.GetByID(uint(id))
	if err != nil {
		h.logger.Error("Ошибка получения patient record", "error", err.Error(), "id", id)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Patient record получен", "id", record.ID)
	c.JSON(200, record)
}

func (h *PatientRecordHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Неверный ID в PatientRecord.Update", "param", idStr)
		c.JSON(400, gin.H{"error": "некорректный ID"})
		return
	}

	var req models.PatientRecordUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "некорректный JSON"})
		return
	}

	if err := h.service.Update(uint(id), &req); err != nil {
		h.logger.Error("Ошибка обновления patient record", "error", err.Error(), "id", id)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Patient record обновлён", "id", id)
	c.JSON(200, gin.H{"message": "запись успешно обновлена"})
}

func (h *PatientRecordHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Неверный ID в PatientRecord.Delete", "param", idStr)
		c.JSON(400, gin.H{"error": "некорректный ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		h.logger.Error("Ошибка удаления patient record", "error", err.Error(), "id", id)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Patient record удалён", "id", id)
	c.JSON(200, gin.H{"message": "запись успешно удалена"})
}
