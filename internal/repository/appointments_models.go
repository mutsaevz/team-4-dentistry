package repository

import (
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	Create(*models.Appointment) error
	Update(*models.Appointment) error
	Delete(uint) error
	GetByID(uint) (*models.Appointment, error)
	Get() ([]models.Appointment, error)
}
type gormAppointmentRepository struct {
	DB *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &gormAppointmentRepository{DB: db}
}

func (r *gormAppointmentRepository) Create(appointment *models.Appointment) error {

	if appointment == nil {
		return constants.Appointments_IS_nil
	}

	dateOnly := time.Date(appointment.StartAt.Year(), appointment.StartAt.Month(), appointment.StartAt.Day(), 0, 0, 0, 0, appointment.StartAt.Location())
	var schedule models.Shedule
	if err := r.DB.Where("doctor_id = ? AND date = ? AND start_time <= ? AND end_time >= ?", appointment.DoctorID, dateOnly, appointment.StartAt, appointment.EndAt).First(&schedule).Error; err != nil {
		
			return constants.ErrTimeNotInSchedule
		}
	


	var count int64
	if err := r.DB.Model(&models.Appointment{}).
		Where("doctor_id = ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.DoctorID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return constants.ErrTimeConflict
	}


	if err := r.DB.Model(&models.Appointment{}).
		Where("patient_id = ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.PatientID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return constants.ErrTimeConflict
	}

	return r.DB.Preload("Service").Preload("Doctor").Preload("Patient").Create(appointment).Error
}

func (r *gormAppointmentRepository) Update(appointment *models.Appointment) error {
	if appointment == nil {
		return constants.Appointments_IS_nil
	}

	dateOnly := time.Date(appointment.StartAt.Year(), appointment.StartAt.Month(), appointment.StartAt.Day(), 0, 0, 0, 0, appointment.StartAt.Location())
	var schedule models.Shedule
	if err := r.DB.Where("doctor_id = ? AND date = ? AND start_time <= ? AND end_time >= ?", appointment.DoctorID, dateOnly, appointment.StartAt, appointment.EndAt).First(&schedule).Error; err != nil {
		
			return constants.ErrTimeNotInSchedule
		}

		var count int64
	if err := r.DB.Model(&models.Appointment{}).
		Where("doctor_id = ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.DoctorID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return constants.ErrTimeConflict
	}

	if err := r.DB.Model(&models.Appointment{}).
		Where("patient_id = ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.PatientID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return constants.ErrTimeConflict
	}

	return r.DB.Preload("Service").Preload("Doctor").Preload("Patient").Save(appointment).Error
}

func (r *gormAppointmentRepository) Delete(id uint) error {

	return r.DB.Delete(&models.Appointment{}, id).Error

}

func (r *gormAppointmentRepository) GetByID(id uint) (*models.Appointment, error) {
	var appointment models.Appointment

	if err := r.DB.Preload("Service").Preload("Doctor").Preload("Patient").First(&appointment, id).Error; err != nil {
		return nil, err
	}

	return &appointment, nil
}

func (r *gormAppointmentRepository) Get() ([]models.Appointment, error) {
	var appointments []models.Appointment

	if err := r.DB.Preload("Service").Preload("Doctor").Preload("Patient").Find(&appointments).Error; err != nil {
		return nil, err
	}

	return appointments, nil
}
