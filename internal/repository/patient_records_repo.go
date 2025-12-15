package repository

import (
	"log/slog"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type PatientRecordRepo interface {
	Create(*models.PatientRecord) error
	GetID(uint) (*models.PatientRecord, error)
	Get() ([]models.PatientRecord, error)
	Update(*models.PatientRecord) error
	Delete(uint) error
}

type gormPatientRecordRepo struct {
	DB *gorm.DB
	logger *slog.Logger
}

func NewPatientRecordRepo(db *gorm.DB, logger *slog.Logger) PatientRecordRepo {
	return &gormPatientRecordRepo{DB: db, logger: logger}
}

func (r *gormPatientRecordRepo) Create(patientRecord *models.PatientRecord) error {
	if patientRecord == nil {
		r.logger.Warn("patientRecord равен nil")
		return constants.PatientRecord_IS_nil
	}

	if err := r.DB.Create(patientRecord).Error; err != nil {
		r.logger.Error("ошибка при создании patientRecord", "ошибка", err)
		return err
	}

	r.logger.Info("patientRecord успешно создан", "patientRecord_id", patientRecord.ID)
	return nil
}

func (r *gormPatientRecordRepo) GetID(ID uint) (*models.PatientRecord, error) {
	r.logger.Debug("получение patientRecord благодаря ID", "patientRecord_id", ID)
	var patientRecord models.PatientRecord

	if err := r.DB.Preload("patient").First(&patientRecord, ID).Error; err != nil {
		r.logger.Error("ошибка при получении patientRecord по ID", "ошибка", err, "patientRecord_id", ID)
		return nil, err
	}

	r.logger.Info("успешное получение patientRecord по ID", "patientRecord_id", ID)

	return &patientRecord, nil
}

func (r *gormPatientRecordRepo) Get() ([]models.PatientRecord, error) {
	r.logger.Warn("Получение всех patient_records")
	var patientRecord []models.PatientRecord

	if err := r.DB.Preload("patient").Find(&patientRecord).Error; err != nil {
		r.logger.Error("ошибка при получении всех patient_records", "ошибка", err)
		return nil, err
	}

	r.logger.Info("успешное получение всех patient_records")
	return patientRecord, nil
}

func (r *gormPatientRecordRepo) Update(patientRecord *models.PatientRecord) error {
	if patientRecord == nil {
		r.logger.Warn("patientRecord равен nil")
		return constants.PatientRecord_IS_nil
	}

	r.logger.Info("Обновление patientRecord", "patientRecord_id", patientRecord.ID)

	if err := r.DB.Save(patientRecord).Error; err != nil {
		r.logger.Error("ошибка при обновлении patientRecord", "ошибка", err, "patientRecord_id", patientRecord.ID)
		return err
	}

	r.logger.Info("успешное обновление patientRecord", "patientRecord_id", patientRecord.ID)
	return nil
}

func (r *gormPatientRecordRepo) Delete(ID uint) error {

	r.logger.Info("Удаление patientRecord по ID", "patientRecord_id", ID)
	
	if err := r.DB.Delete(&models.PatientRecord{}, ID).Error; err != nil {
		r.logger.Error("ошибка при удалении patientRecord", "ошибка", err, "patientRecord_id", ID)
		return err
	}

	r.logger.Info("успешное удаление patientRecord", "patientRecord_id", ID)
	return nil
}
