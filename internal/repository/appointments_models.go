package repository

import (
	"errors"
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	Delete(uint) error
	GetByID(uint) (*models.Appointment, error)
	Get() ([]models.Appointment, error)
	Transaction(func(tx *gorm.DB) error) error
	CreateTx(tx *gorm.DB, appointment *models.Appointment) error
	UpdateTx(tx *gorm.DB, appointment *models.Appointment) error
	GetByPatientID(patientID uint)([]models.Appointment, error)
}
type gormAppointmentRepository struct {
	DB *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &gormAppointmentRepository{DB: db}
}

func (r *gormAppointmentRepository) Delete(id uint) error {

	return r.DB.Delete(&models.Appointment{}, id).Error

}

func (r *gormAppointmentRepository) GetByID(id uint) (*models.Appointment, error) {
	var appointment models.Appointment

	if err := r.DB.First(&appointment, id).Error; err != nil {
		return nil, err
	}

	return &appointment, nil
}

func (r *gormAppointmentRepository) Get() ([]models.Appointment, error) {
	var appointments []models.Appointment

	if err := r.DB.Find(&appointments).Error; err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *gormAppointmentRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.DB.Transaction(fn)
}

func (r *gormAppointmentRepository) CreateTx(tx *gorm.DB, appointment *models.Appointment) error {
	if appointment == nil {
		return constants.Appointments_IS_nil
	}

	dateOnly := time.Date(appointment.StartAt.Year(), appointment.StartAt.Month(), appointment.StartAt.Day(), 0, 0, 0, 0, appointment.StartAt.Location())
	var schedule models.Schedule
	if err := tx.Where("doctor_id = ? AND date = ? AND start_time <= ? AND end_time >= ?", appointment.DoctorID, dateOnly, appointment.StartAt, appointment.EndAt).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrTimeNotInSchedule
		}
		return err
	}

	var count int64
	if err := tx.Model(&models.Appointment{}).
		Where("doctor_id = ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.DoctorID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return constants.ErrTimeConflict
	}

	if err := tx.Model(&models.Appointment{}).
		Where("patient_id = ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.PatientID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return constants.ErrTimeConflict
	}

	return tx.Create(appointment).Error
}

func (r *gormAppointmentRepository) UpdateTx(tx *gorm.DB, appointment *models.Appointment) error {
	if appointment == nil {
		return constants.Appointments_IS_nil
	}

	dateOnly := time.Date(appointment.StartAt.Year(), appointment.StartAt.Month(), appointment.StartAt.Day(), 0, 0, 0, 0, appointment.StartAt.Location())
	var schedule models.Schedule
	if err := tx.Where("doctor_id = ? AND date = ? AND start_time <= ? AND end_time >= ?", appointment.DoctorID, dateOnly, appointment.StartAt, appointment.EndAt).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrTimeNotInSchedule
		}
		return err
	}

	var count int64
	if err := tx.Model(&models.Appointment{}).
		Where("doctor_id = ? AND id <> ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.DoctorID, appointment.ID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return constants.ErrTimeConflict
	}

	if err := tx.Model(&models.Appointment{}).
		Where("patient_id = ? AND id <> ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.PatientID, appointment.ID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return constants.ErrTimeConflict
	}

	return tx.Save(appointment).Error
}

func (r *gormAppointmentRepository) GetByPatientID(patientID uint)([]models.Appointment, error){
	var appointment []models.Appointment

	if err := r.DB.Where("user_id = ?", patientID).Find(&appointment); err != nil {
		return nil, constants.User_appointments_Not_Found
	}

	return appointment, nil
}
