package services

import (
	"errors"
	"log/slog"
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
	logger            *slog.Logger
}

func NewAppointmentService(service repository.ServiceRepository, appointments repository.AppointmentRepository, logger *slog.Logger) AppointmentService {
	return &appointmentService{serviceRepository: service, appointments: appointments, logger: logger}
}

func (r *appointmentService) Create(req *models.AppointmentCreateRequest) (*models.Appointment, error) {
	r.logger.Debug("создание appointment вызвано", "doctor_id", req.DoctorID, "patient_id", req.PatientID, "service_id", req.ServiceID)

	if req == nil {
		r.logger.Warn("получен nil AppointmentCreateRequest")
		return nil, constants.AppointmentCreateRequest_IS_nil
	}
	if err := r.validate(req); err != nil {
		r.logger.Warn("валидация создания appointment провалилась", "error", err)
		return nil, err
	}

	service, err := r.serviceRepository.GetByID(req.ServiceID)
	if err != nil {
		r.logger.Error("не удалось получить service для appointment", "error", err, "service_id", req.ServiceID)
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
			r.logger.Error("ошибка при создании appointment в транзакции", "error", err)
			return err
		}

		if err := tx.
			Model(&models.Schedule{}).
			Where("doctor_id = ? AND start_time = ?", req.DoctorID, req.StartAt).
			Update("is_available", false).Error; err != nil {
			r.logger.Error("ошибка при обновлении доступности расписания", "error", err)
			return err
		}

		return nil
	})

	if err != nil {
		r.logger.Error("транзакция создания appointment провалилась", "error", err)
		return nil, err
	}
	r.logger.Info("appointment создан", "appointment_id", appointment.ID)
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
	r.logger.Debug("обновление appointment вызвано", "appointment_id", id)

	appointments, err := r.appointments.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("appointment не найден для обновления", "appointment_id", id)
			return constants.ErrInvalidAppointmentID
		}
		r.logger.Error("ошибка при получении appointment для обновления", "error", err, "appointment_id", id)
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
		r.logger.Error("транзакция обновления appointment провалилась", "error", err, "appointment_id", id)
		return constants.ErrUpdateAppointments
	}

	r.logger.Info("appointment успешно обновлён", "appointment_id", id)
	return nil
}

func (r *appointmentService) Delete(id uint) error {
	r.logger.Debug("удаление appointment вызвано", "appointment_id", id)
	if id <= 0 {
		r.logger.Warn("некорректный id для удаления appointment", "appointment_id", id)
		return constants.ErrInvalidAppointmentID
	}

	if err := r.appointments.Delete(id); err != nil {
		r.logger.Error("ошибка при удалении appointment", "error", err, "appointment_id", id)
		return constants.ErrDeleteAppointments
	}

	r.logger.Info("appointment удалён", "appointment_id", id)
	return nil
}

func (r *appointmentService) GetByID(id uint) (*models.Appointment, error) {
	r.logger.Debug("получение appointment по ID вызвано", "appointment_id", id)
	if id <= 0 {
		r.logger.Warn("некорректный id для GetByID appointment", "appointment_id", id)
		return nil, constants.ErrInvalidAppointmentID
	}

	appointment, err := r.appointments.GetByID(id)
	if err != nil {
		r.logger.Error("ошибка при получении appointment по id", "error", err, "appointment_id", id)
		return nil, constants.ErrGetByIDAppointments
	}
	r.logger.Info("appointment получен по id", "appointment_id", id)
	return appointment, nil
}

func (r *appointmentService) GetAll() ([]models.Appointment, error) {
	r.logger.Debug("получение всех appointments вызвано")
	appointments, err := r.appointments.Get()
	if err != nil {
		r.logger.Error("ошибка при получении всех appointments", "error", err)
		return nil, constants.ErrGetAppointments
	}
	r.logger.Info("appointments получены", "count", len(appointments))
	return appointments, nil
}

func (r *appointmentService) GetByPatientID(patientID uint) ([]models.Appointment, error) {
	r.logger.Debug("GetByPatientID вызван", "patient_id", patientID)
	if patientID <= 0 {
		r.logger.Warn("некорректный patient_id для GetByPatientID", "patient_id", patientID)
		return nil, constants.PatientIDIsIncorrect
	}

	appointments, err := r.appointments.GetByPatientID(patientID)
	if err != nil {
		r.logger.Error("ошибка при получении appointments по patient_id", "error", err, "patient_id", patientID)
		return nil, err
	}

	r.logger.Info("appointments получены по patient_id", "patient_id", patientID, "count", len(appointments))
	return appointments, nil
}
