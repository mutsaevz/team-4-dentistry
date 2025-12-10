package services

import (
	"context"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, req models.SheduleCreateRequest) (*models.Shedule, error)

	GetScheduleByID(ctx context.Context, id uint) (*models.Shedule, error)

	ListSchedules(ctx context.Context) ([]models.Shedule, error)

	UpdateSchedule(ctx context.Context, id uint, req models.SheduleUpdateRequest) (*models.Shedule, error)

	DeleteSchedule(ctx context.Context, id uint) error
}

type scheduleService struct {
	schedule repository.SheduleRepositroy
	doctor   repository.DoctorRepository
}

func NewSheduleService(repoShedule repository.SheduleRepositroy, repoDoctor repository.DoctorRepository) ScheduleService {
	return &scheduleService{
		schedule: repoShedule,
		doctor:   repoDoctor,
	}
}

func (s *scheduleService) CreateSchedule(ctx context.Context, req models.SheduleCreateRequest) (*models.Shedule, error) {
	if err := s.ValidateScheduleCreate(req); err != nil {
		return nil, err
	}

	schedule := &models.Shedule{
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

func (s *scheduleService) GetScheduleByID(ctx context.Context, id uint) (*models.Shedule, error) {
	return s.schedule.GetByID(id, ctx)
}

func (s *scheduleService) ListSchedules(ctx context.Context) ([]models.Shedule, error) {
	return s.schedule.GetAll(ctx)
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, id uint, req models.SheduleUpdateRequest) (*models.Shedule, error) {
	schedule, err := s.schedule.GetByID(id, ctx)
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

func (s *scheduleService) ValidateScheduleCreate(req models.SheduleCreateRequest) error {
	if req.DoctorID == 0 {
		return constants.ErrInvalidDoctorID
	}

	if req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime) {
		return constants.ErrInvalidTimeRange
	}

	if req.RoomNumber == "" {
		return constants.ErrInvalidRoomNumber
	}

	return nil
}

func (s *scheduleService) ValidateScheduleUpdate(existing *models.Shedule, req models.SheduleUpdateRequest) error {
	if req.DoctorID != nil && *req.DoctorID != 0 {
		existing.DoctorID = *req.DoctorID
	}

	if req.StartTime != nil {
		existing.StartTime = *req.StartTime
	}

	if req.EndTime != nil {
		existing.EndTime = *req.EndTime
	}

	if req.RoomNumber != nil && *req.RoomNumber != "" {
		existing.RoomNumber = *req.RoomNumber
	}

	if req.IsAvailable != nil {
		existing.IsAvailable = *req.IsAvailable
	}

	return nil
}
