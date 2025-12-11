package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type ScheduleRepositroy interface {
	Create(context.Context, *models.Schedule) error

	GetAll(context.Context) ([]models.Schedule, error)

	GetByDoctorID(uint, context.Context) (*models.Schedule, error)

	GetByDateRange(context.Context, uint, time.Time, time.Time) (*models.Schedule, error)

	CheckConflict(context.Context, uint, time.Time, time.Time) (bool, error)

	Update(context.Context, *models.Schedule) error

	Delete(context.Context, uint) error

	DeleteByDoctorID(context.Context, uint) error

	GetAvailableSlots(context.Context, uint) ([]models.Schedule, error)
}

type gormScheduleRepository struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepositroy {
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

func (r *gormScheduleRepository) GetByDoctorID(id uint, ctx context.Context) (*models.Schedule, error) {
	var schedule models.Schedule

	if err := r.DB.WithContext(ctx).Where("doctor_id = ?", id).First(&schedule).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *gormScheduleRepository) GetByDateRange(ctx context.Context, doctorID uint, start, end time.Time) (*models.Schedule, error) {
	var schedule models.Schedule

	if err := r.DB.WithContext(ctx).Where("doctor_id = ? AND date >= ? AND date <= ?", doctorID, start, end).First(&schedule).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *gormScheduleRepository) CheckConflict(ctx context.Context, doctorID uint, start, end time.Time) (bool, error) {
	var count int64

	err := r.DB.WithContext(ctx).Model(&models.Schedule{}).
		Where("doctor_id = ? AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?))",
			doctorID, end, end, start, start, start, end).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
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

func (r *gormScheduleRepository) GetAvailableSlots(ctx context.Context, doctorID uint) ([]models.Schedule, error) {
	var schedules []models.Schedule

	if err := r.DB.WithContext(ctx).Where("doctor_id = ? AND date >= ?", doctorID, time.Now()).Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}
