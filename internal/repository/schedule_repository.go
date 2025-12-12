package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type ScheduleRepository interface {
	Create(context.Context, *models.Schedule) error

	GetAll(context.Context) ([]models.Schedule, error)

	GetByID(context.Context, uint) (*models.Schedule, error)

	GetSchedulesByDoctorID(context.Context, uint) ([]models.Schedule, error)

	Update(context.Context, *models.Schedule) error

	Delete(context.Context, uint) error

	DeleteByDoctorID(context.Context, uint) error

	GetAvailableSlots(context.Context, uint, time.Time, time.Time) ([]models.Schedule, error)
}

type gormScheduleRepository struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &gormScheduleRepository{DB: db}
}

func (r *gormScheduleRepository) Create(ctx context.Context, schedule *models.Schedule) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}

	return r.DB.WithContext(ctx).Create(schedule).Error
}

func (r *gormScheduleRepository) GetAll(ctx context.Context) ([]models.Schedule, error) {
	var schedules []models.Schedule

	if err := r.DB.WithContext(ctx).Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *gormScheduleRepository) GetByID(ctx context.Context, id uint) (*models.Schedule, error) {
	var schedule models.Schedule

	if err := r.DB.WithContext(ctx).First(&schedule, id).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *gormScheduleRepository) GetSchedulesByDoctorID(ctx context.Context, doctorID uint) ([]models.Schedule, error) {
	var schedule []models.Schedule

	if err := r.DB.WithContext(ctx).Where("doctor_id = ?", doctorID).Find(&schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *gormScheduleRepository) Update(ctx context.Context, schedule *models.Schedule) error {
	if schedule == nil {
		return nil
	}

	return r.DB.WithContext(ctx).Save(schedule).Error
}

func (r *gormScheduleRepository) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&models.Schedule{}, id).Error
}

func (r *gormScheduleRepository) DeleteByDoctorID(ctx context.Context, doctorID uint) error {
	return r.DB.WithContext(ctx).Where("doctor_id = ?", doctorID).Delete(&models.Schedule{}).Error
}

func (r *gormScheduleRepository) GetAvailableSlots(
	ctx context.Context,
	doctorID uint,
	start, end time.Time,
) ([]models.Schedule, error) {

	start = start.Truncate(24 * time.Hour)
	end = start.AddDate(0, 0, 7)

	var schedules []models.Schedule

	err := r.DB.WithContext(ctx).
		Where("doctor_id = ?", doctorID).
		Where("is_available = ?", true).
		Where("start_time >= ? AND end_time <= ?", start, end).
		Order("start_time ASC").
		Find(&schedules).Error

	if err != nil {
		return nil, err
	}

	return schedules, nil
}
