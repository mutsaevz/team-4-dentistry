package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
)

type DoctorService interface {
	CreateDoctor(context.Context, models.DoctorCreateRequest) (*models.Doctor, error)

	GetDoctorByID(context.Context, uint) (*models.Doctor, error)

	ListDoctors(context.Context, models.DoctorQueryParams) ([]models.Doctor, error)

	UpdateDoctor(context.Context, uint, models.DoctorUpdateRequest) (*models.Doctor, error)

	DeleteDoctor(context.Context, uint) error

	GetDoctorServices(context.Context, uint) ([]models.Service, error)

	GetScheduleByDoctorID(context.Context, uint) ([]models.Schedule, error)
}

type doctorService struct {
	doctors  repository.DoctorRepository
	service  repository.ServiceRepository
	schedule repository.ScheduleRepository
	logger   *slog.Logger
}

func NewDoctorService(
	doctors repository.DoctorRepository,
	service repository.ServiceRepository,
	schedule repository.ScheduleRepository,
	logger *slog.Logger,
) DoctorService {
	return &doctorService{
		doctors:  doctors,
		service:  service,
		schedule: schedule,
		logger:   logger,
	}
}

func (s *doctorService) CreateDoctor(ctx context.Context, req models.DoctorCreateRequest) (*models.Doctor, error) {
	s.logger.Debug("CreateDoctor вызван", "user_id", req.UserID)

	if err := s.ValidateCreateDoctor(req); err != nil {
		s.logger.Error("валидация CreateDoctor провалилась", "error", err, "user_id", req.UserID)
		return nil, err
	}

	doctor := &models.Doctor{
		UserID:          req.UserID,
		Specialization:  req.Specialization,
		ExperienceYears: req.ExperienceYears,
		Bio:             req.Bio,
		AvgRating:       0,
		RoomNumber:      req.RoomNumber,
	}

	if err := s.doctors.Create(ctx, doctor); err != nil {
		s.logger.Error("ошибка при создании врача в репозитории", "error", err, "user_id", req.UserID)
		return nil, err
	}
	s.logger.Info("врач создан", "doctor_id", doctor.ID, "user_id", req.UserID)
	return doctor, nil
}

func (s *doctorService) GetDoctorByID(ctx context.Context, id uint) (*models.Doctor, error) {
	s.logger.Debug("получение врача по ID вызвано", "doctor_id", id)
	doctor, err := s.doctors.GetByID(id, ctx)
	if err != nil {
		s.logger.Error("ошибка при получении врача по ID", "error", err, "doctor_id", id)
		return nil, err
	}
	s.logger.Info("врач получен по ID", "doctor_id", id)
	return doctor, nil
}

func (s *doctorService) ListDoctors(ctx context.Context, params models.DoctorQueryParams) ([]models.Doctor, error) {
	s.logger.Debug("получение списка врачей вызвано", "params", params)
	doctors, err := s.doctors.GetAll(params, ctx)
	if err != nil {
		s.logger.Error("ошибка при получении списка врачей", "error", err)
		return nil, err
	}
	s.logger.Info("список врачей получен", "count", len(doctors))
	return doctors, nil
}

func (s *doctorService) GetDoctorServices(ctx context.Context, doctorID uint) ([]models.Service, error) {
	s.logger.Debug("получение услуг врача вызвано", "doctor_id", doctorID)
	svcs, err := s.service.GetServicesByDoctorID(ctx, doctorID)
	if err != nil {
		s.logger.Error("ошибка при получении услуг врача", "error", err, "doctor_id", doctorID)
		return nil, err
	}
	s.logger.Info("услуги врача получены", "doctor_id", doctorID, "count", len(svcs))
	return svcs, nil
}

func (s *doctorService) GetScheduleByDoctorID(ctx context.Context, doctorID uint) ([]models.Schedule, error) {
	s.logger.Debug("получение расписания врача вызвано", "doctor_id", doctorID)
	sch, err := s.schedule.GetSchedulesByDoctorID(ctx, doctorID)
	if err != nil {
		s.logger.Error("ошибка при получении расписания врача", "error", err, "doctor_id", doctorID)
		return nil, err
	}
	s.logger.Info("расписание врача получено", "doctor_id", doctorID, "count", len(sch))
	return sch, nil
}

func (s *doctorService) UpdateDoctor(ctx context.Context, id uint, req models.DoctorUpdateRequest) (*models.Doctor, error) {
	s.logger.Debug("обновление врача вызвано", "doctor_id", id)
	doctor, err := s.doctors.GetByID(id, ctx)
	if err != nil {
		s.logger.Error("ошибка при получении врача для обновления", "error", err, "doctor_id", id)
		return nil, err
	}

	if *req.Specialization != "" {
		doctor.Specialization = *req.Specialization
	}
	if req.ExperienceYears != nil {
		doctor.ExperienceYears = *req.ExperienceYears
	}
	if req.Bio != nil {
		doctor.Bio = *req.Bio
	}
	if req.RoomNumber != nil {
		doctor.RoomNumber = *req.RoomNumber
	}

	s.logger.Info("врач обновлён", "doctor_id", id)
	return doctor, nil
}

func (s *doctorService) DeleteDoctor(ctx context.Context, id uint) error {
	s.logger.Debug("удаление врача вызвано", "doctor_id", id)
	if err := s.doctors.Delete(ctx, id); err != nil {
		s.logger.Error("ошибка при удалении врача", "error", err, "doctor_id", id)
		return err
	}
	s.logger.Info("врач удалён", "doctor_id", id)
	return nil
}

func (s *doctorService) ValidateCreateDoctor(req models.DoctorCreateRequest) error {
	if req.UserID <= 0 {
		return errors.New("")
	}
	if req.Specialization == "" {
		return errors.New("")
	}
	if req.RoomNumber <= 0 {
		return errors.New("")
	}
	if req.ExperienceYears < 0 {
		return errors.New("")
	}
	if req.Bio == "" {
		return errors.New("")
	}

	return nil
}
