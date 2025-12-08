package repository

import (

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

	return r.DB.Create(appointment).Error
}

func (r *gormAppointmentRepository) Update(appointment *models.Appointment) error {
	if appointment == nil {
		return constants.Appointments_IS_nil
	}

	return r.DB.Save(appointment).Error
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
