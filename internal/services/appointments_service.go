package services

import (
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
)

type AppointmentService interface {
}

type appointmentService struct {
	serviceRepository repository.ServiceRepository
	appointments repository.AppointmentRepository
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
	
	if err := r.appointments.Create(appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

func (r *appointmentService)validate(req *models.AppointmentCreateRequest) error {
	if req.DoctorID <= 0  {
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

func (r *appointmentService) Update(id uint, req models.AppointmentUpdateRequest)(error){

	appointments, err := r.appointments.GetByID(id) 

	if err != nil{
		return constants.ErrInvalidAppointmentID
	}

	if req.DoctorID != nil && *req.DoctorID <= 0 {
		return constants.DoctorIDIsIncorrect
	}

	if req.PatientID != nil && *req.PatientID <= 0 {
		return constants.PatientIDIsIncorrect
	}

	if req.ServiceID != nil && *req.ServiceID <= 0{
		return constants.ServiceIDIsIncorrect
	}

	if req.StartAt != nil && req.StartAt.Before(time.Now()){
		return constants.ErrInvalidAppointmentTime
	}

	if req.Price != nil && *req.Price <= 0 {
		return constants.ErrInvalidPrice
	}

	if err := r.appointments.Update(appointments); err != nil {
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