package services

import (
	"context"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, req models.ScheduleCreateRequest) (*models.Schedule, error)

	GetScheduleByID(ctx context.Context, id uint) (*models.Schedule, error)

	ListSchedules(ctx context.Context) ([]models.Schedule, error)

	UpdateSchedule(ctx context.Context, id uint, req models.ScheduleUpdateRequest) (*models.Schedule, error)

	DeleteSchedule(ctx context.Context, id uint) error
}

type scheduleService struct {
	schedule repository.ScheduleRepository
	doctor   repository.DoctorRepository
}

func NewScheduleService(repoSchedule repository.ScheduleRepository, repoDoctor repository.DoctorRepository) ScheduleService {
	return &scheduleService{
		schedule: repoSchedule,
		doctor:   repoDoctor,
	}
}

func (s *scheduleService) CreateSchedule(ctx context.Context, req models.ScheduleCreateRequest) (*models.Schedule, error) {
	if err := s.ValidateScheduleCreate(req); err != nil {
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
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleService) GetScheduleByID(ctx context.Context, id uint) (*models.Schedule, error) {
	return s.schedule.GetByDoctorID(id, ctx)
}

func (s *scheduleService) ListSchedules(ctx context.Context) ([]models.Schedule, error) {
	return s.schedule.GetAll(ctx)
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, id uint, req models.ScheduleUpdateRequest) (*models.Schedule, error) {
	schedule, err := s.schedule.GetByDoctorID(id, ctx)
	if err != nil {
		return nil, err
	}

	if err := s.ValidateScheduleUpdate(schedule, req); err != nil {
		return nil, err
	}

	if err := s.schedule.Update(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, id uint) error {
	return s.schedule.Delete(ctx, id)
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
