package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, req models.ScheduleCreateRequest) (*models.Schedule, error)

	GetSchedulesByID(ctx context.Context, id uint) ([]models.Schedule, error)

	ListSchedules(ctx context.Context) ([]models.Schedule, error)

	UpdateSchedule(ctx context.Context, id uint, req models.ScheduleUpdateRequest) (*models.Schedule, error)

	DeleteSchedule(ctx context.Context, id uint) error

	GetAvailableSlots(ctx context.Context, doctorID uint, week int) ([]models.Schedule, error)
}

type scheduleService struct {
	schedule repository.ScheduleRepository
	doctor   repository.DoctorRepository
	logger   *slog.Logger
}

func NewScheduleService(repoSchedule repository.ScheduleRepository, repoDoctor repository.DoctorRepository, logger *slog.Logger) ScheduleService {
	return &scheduleService{
		schedule: repoSchedule,
		doctor:   repoDoctor,
		logger:   logger,
	}
}

func (s *scheduleService) CreateSchedule(ctx context.Context, req models.ScheduleCreateRequest) (*models.Schedule, error) {
	s.logger.Debug("CreateSchedule вызван", "doctor_id", req.DoctorID, "date", req.Date)
	if err := s.ValidateScheduleCreate(req); err != nil {
		s.logger.Error("валидация CreateSchedule провалилась", "error", err)
		return nil, err
	}

	schedule := &models.Schedule{
		DoctorID:    req.DoctorID,
		Date:        req.Date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		RoomNumber:  req.RoomNumber,
		IsAvailable: true,
	}

	if err := s.schedule.Create(ctx, schedule); err != nil {
		s.logger.Error("ошибка при создании schedule", "error", err)
		return nil, err
	}

	s.logger.Info("schedule создан", "schedule_id", schedule.ID, "doctor_id", schedule.DoctorID)
	return schedule, nil
}

func (s *scheduleService) GetSchedulesByID(ctx context.Context, id uint) ([]models.Schedule, error) {
	s.logger.Debug("GetSchedulesByID вызван", "doctor_id", id)
	sch, err := s.schedule.GetSchedulesByDoctorID(ctx, id)
	if err != nil {
		s.logger.Error("ошибка при получении расписания по doctor_id", "error", err, "doctor_id", id)
		return nil, err
	}
	s.logger.Info("расписание получено по doctor_id", "doctor_id", id, "count", len(sch))
	return sch, nil
}

func (s *scheduleService) ListSchedules(ctx context.Context) ([]models.Schedule, error) {
	s.logger.Debug("ListSchedules вызван")
	sch, err := s.schedule.GetAll(ctx)
	if err != nil {
		s.logger.Error("ошибка при получении всех расписаний", "error", err)
		return nil, err
	}
	s.logger.Info("все расписания получены", "count", len(sch))
	return sch, nil
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, id uint, req models.ScheduleUpdateRequest) (*models.Schedule, error) {
	s.logger.Debug("UpdateSchedule вызван", "schedule_id", id)
	schedule, err := s.schedule.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("ошибка при получении schedule для обновления", "error", err, "schedule_id", id)
		return nil, err
	}

	if err := s.ValidateScheduleUpdate(schedule, req); err != nil {
		s.logger.Error("валидация UpdateSchedule провалилась", "error", err, "schedule_id", id)
		return nil, err
	}

	if err := s.schedule.Update(ctx, schedule); err != nil {
		s.logger.Error("ошибка при обновлении schedule", "error", err, "schedule_id", id)
		return nil, err
	}

	s.logger.Info("schedule успешно обновлен", "schedule_id", id)
	return schedule, nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, id uint) error {
	s.logger.Debug("DeleteSchedule вызван", "schedule_id", id)
	if err := s.schedule.Delete(ctx, id); err != nil {
		s.logger.Error("ошибка при удалении schedule", "error", err, "schedule_id", id)
		return err
	}
	s.logger.Info("schedule удален", "schedule_id", id)
	return nil
}

func (s *scheduleService) ValidateScheduleCreate(req models.ScheduleCreateRequest) error {
	if req.DoctorID <= 0 {
		return constants.ErrInvalidDoctorID
	}

	if req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime) {
		return constants.ErrInvalidTimeRange
	}

	if req.RoomNumber <= 0 {
		return constants.ErrInvalidRoomNumber
	}

	return nil
}

func (s *scheduleService) ValidateScheduleUpdate(existing *models.Schedule, req models.ScheduleUpdateRequest) error {
	if req.DoctorID != nil && *req.DoctorID != 0 {
		existing.DoctorID = *req.DoctorID
	}

	if req.StartTime != nil {
		existing.StartTime = *req.StartTime
	}

	if req.EndTime != nil {
		existing.EndTime = *req.EndTime
	}

	if req.RoomNumber != nil && *req.RoomNumber != 0 {
		existing.RoomNumber = *req.RoomNumber
	}

	if req.IsAvailable != nil {
		existing.IsAvailable = *req.IsAvailable
	}

	return nil
}

func (s *scheduleService) GetAvailableSlots(ctx context.Context, doctorID uint, week int) ([]models.Schedule, error) {
	start := time.Now().AddDate(0, 0, 7*week)
	s.logger.Debug("GetAvailableSlots вызван", "doctor_id", doctorID, "week", week, "start", start)
	slots, err := s.schedule.GetAvailableSlots(ctx, doctorID, start)
	if err != nil {
		s.logger.Error("ошибка при получении доступных слотов", "error", err, "doctor_id", doctorID)
		return nil, err
	}
	s.logger.Info("доступные слоты получены", "doctor_id", doctorID, "count", len(slots))
	return slots, nil
}
