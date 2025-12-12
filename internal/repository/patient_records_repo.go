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
		return constants.PatientRecord_IS_nil
	}

	return r.DB.Create(patientRecord).Error
}

func (r *gormPatientRecordRepo) GetID(ID uint) (*models.PatientRecord, error) {
	var patientRecord models.PatientRecord

	if err := r.DB.Preload("patient").First(&patientRecord, ID).Error; err != nil {
		return nil, err
	}

	return &patientRecord, nil
}

func (r *gormPatientRecordRepo) Get() ([]models.PatientRecord, error) {
	var patientRecord []models.PatientRecord

	if err := r.DB.Preload("patient").Find(&patientRecord).Error; err != nil {
		return nil, err
	}

	return patientRecord, nil
}

func (r *gormPatientRecordRepo) Update(patientRecord *models.PatientRecord) error {
	if patientRecord == nil {
		return constants.PatientRecord_IS_nil
	}
	
	return r.DB.Save(patientRecord).Error
}

func (r *gormPatientRecordRepo) Delete(ID uint) error {
	return r.DB.Delete(&models.PatientRecord{}, ID).Error
}
