package services

import (
	"errors"
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
	"gorm.io/gorm"
)

type AppointmentService interface {
	Create(req *models.AppointmentCreateRequest) (*models.Appointment, error)
	Update(id uint, req *models.AppointmentUpdateRequest) error
	Delete(id uint) error
	GetByID(id uint) (*models.Appointment, error)
	GetAll() ([]models.Appointment, error)
	GetByPatientID(patientID uint) ([]models.Appointment, error)
}

type appointmentService struct {
	serviceRepository repository.ServiceRepository
	appointments      repository.AppointmentRepository
}

func NewAppointmentService(service repository.ServiceRepository, appointments repository.AppointmentRepository) AppointmentService {
	return &appointmentService{serviceRepository: service, appointments: appointments}
}

func (r *appointmentService) Create(req *models.AppointmentCreateRequest) (*models.Appointment, error) {
	if req == nil {
		return nil, constants.AppointmentCreateRequest_IS_nil
	}

	if err := r.validate(req); err != nil {
		return nil, err
	}

	service, err := r.serviceRepository.GetByID(req.ServiceID)
	if err != nil {
		return nil, err
	}

	duration := service.Duration

	appointment := &models.Appointment{
		PatientID: req.PatientID,
		DoctorID:  req.DoctorID,
		ServiceID: req.ServiceID,
		StartAt:   req.StartAt,
		EndAt:     req.StartAt.Add(time.Duration(duration) * time.Minute),
		Price:     req.Price,
	}

	err = r.appointments.Transaction(func(tx *gorm.DB) error {

		if err := r.appointments.CreateTx(tx, appointment); err != nil {
			return err
		}

		if err := tx.
			Model(&models.Schedule{}).
			Where("doctor_id = ? AND start_time = ?", req.DoctorID, req.StartAt).
			Update("is_available", false).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return appointment, nil
}

func (r *appointmentService) validate(req *models.AppointmentCreateRequest) error {
	if req.DoctorID <= 0 {
		return constants.DoctorIDIsIncorrect
	}

	if req.PatientID <= 0 {
		return constants.PatientIDIsIncorrect
	}

	if req.ServiceID <= 0 {
		return constants.ServiceIDIsIncorrect
	}

	if req.StartAt.Before(time.Now()) {
		return constants.ErrInvalidAppointmentTime
	}

	if req.Price < 0 {
		return constants.ErrInvalidPrice
	}

	return nil
}

func (r *appointmentService) Update(id uint, req *models.AppointmentUpdateRequest) error {

	appointments, err := r.appointments.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrInvalidAppointmentID
		}
		return err
	}

	if req.DoctorID != nil && *req.DoctorID <= 0 {
		return constants.DoctorIDIsIncorrect
	}

	if req.PatientID != nil && *req.PatientID <= 0 {
		return constants.PatientIDIsIncorrect
	}

	if req.ServiceID != nil && *req.ServiceID <= 0 {
		return constants.ServiceIDIsIncorrect
	}

	if req.StartAt != nil && req.StartAt.Before(time.Now()) {
		return constants.ErrInvalidAppointmentTime
	}

	if req.Price != nil && *req.Price < 0 {
		return constants.ErrInvalidPrice
	}

	if req.DoctorID != nil {
		appointments.DoctorID = *req.DoctorID
	}

	if req.PatientID != nil {
		appointments.PatientID = *req.PatientID
	}

	if req.ServiceID != nil {
		appointments.ServiceID = *req.ServiceID
	}

	if req.StartAt != nil {
		appointments.StartAt = *req.StartAt

		service, err := r.serviceRepository.GetByID(appointments.ServiceID)
		if err != nil {
			return err
		}
		duration := service.Duration
		appointments.EndAt = appointments.StartAt.Add(time.Duration(duration) * time.Minute)
	}

	if req.Price != nil {
		appointments.Price = *req.Price
	}

	if err := r.appointments.Transaction(func(tx *gorm.DB) error {
		return r.appointments.UpdateTx(tx, appointments)
	}); err != nil {
		return constants.ErrUpdateAppointments
	}

	return nil
}

func (r *appointmentService) Delete(id uint) error {
	if id <= 0 {
		return constants.ErrInvalidAppointmentID
	}

	if err := r.appointments.Delete(id); err != nil {
		return constants.ErrDeleteAppointments
	}

	return nil
}

func (r *appointmentService) GetByID(id uint) (*models.Appointment, error) {
	if id <= 0 {
		return nil, constants.ErrInvalidAppointmentID
	}

	appointment, err := r.appointments.GetByID(id)
	if err != nil {
		return nil, constants.ErrGetByIDAppointments
	}
	return appointment, nil
}

func (r *appointmentService) GetAll() ([]models.Appointment, error) {
	appointments, err := r.appointments.Get()

	if err != nil {
		return nil, constants.ErrGetAppointments
	}

	return appointments, nil
}

func (r *appointmentService) GetByPatientID(patientID uint) ([]models.Appointment, error) {

	if patientID <= 0 {
		return nil, constants.PatientIDIsIncorrect
	}

	appointments, err := r.appointments.GetByPatientID(patientID)
	if err != nil {
		return nil, err
	}

	return appointments, nil
}
