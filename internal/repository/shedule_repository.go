package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type SheduleRepositroy interface {
	Create(context.Context, *models.Shedule) error

	GetAll(context.Context) ([]models.Shedule, error)

	GetByDoctorID(uint, context.Context) (*models.Shedule, error)

	GetByDateRange(context.Context, uint, time.Time, time.Time) (*models.Shedule, error)

	CheckConflict(context.Context, uint, time.Time, time.Time) (bool, error)

	Update(context.Context, *models.Shedule) error

	Delete(context.Context, uint) error

	DeleteByDoctorID(context.Context, uint) error

	GetAvailableSlots(context.Context, uint) ([]models.Shedule, error)
}

type gormSheduleRepository struct {
	DB *gorm.DB
}

func NewSheduleRepository(db *gorm.DB) SheduleRepositroy {
	return &gormSheduleRepository{DB: db}
}

func (r *gormSheduleRepository) Create(ctx context.Context, schedule *models.Shedule) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}

	return r.DB.WithContext(ctx).Create(schedule).Error
}

func (r *gormSheduleRepository) GetAll(ctx context.Context) ([]models.Shedule, error) {
	var schedules []models.Shedule

	if err := r.DB.WithContext(ctx).Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *gormSheduleRepository) GetByDoctorID(id uint, ctx context.Context) (*models.Shedule, error) {
	var schedule models.Shedule

	if err := r.DB.WithContext(ctx).Where("doctor_id = ?", id).First(&schedule).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *gormSheduleRepository) GetByDateRange(ctx context.Context, doctorID uint, start, end time.Time) (*models.Shedule, error) {
	var schedule models.Shedule

	if err := r.DB.WithContext(ctx).Where("doctor_id = ? AND date >= ? AND date <= ?", doctorID, start, end).First(&schedule).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *gormSheduleRepository) CheckConflict(ctx context.Context, doctorID uint, start, end time.Time) (bool, error) {
	var count int64

	err := r.DB.WithContext(ctx).Model(&models.Shedule{}).
		Where("doctor_id = ? AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?))",
			doctorID, end, end, start, start, start, end).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *gormSheduleRepository) Update(ctx context.Context, schedule *models.Shedule) error {
	if schedule == nil {
		return nil
	}

	return r.DB.WithContext(ctx).Save(schedule).Error
}

func (r *gormSheduleRepository) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&models.Shedule{}, id).Error
}

func (r *gormSheduleRepository) DeleteByDoctorID(ctx context.Context, doctorID uint) error {
	return r.DB.WithContext(ctx).Where("doctor_id = ?", doctorID).Delete(&models.Shedule{}).Error
}

func (r *gormSheduleRepository) GetAvailableSlots(ctx context.Context, doctorID uint) ([]models.Shedule, error) {
	var schedules []models.Shedule

	if err := r.DB.WithContext(ctx).Where("doctor_id = ? AND date >= ?", doctorID, time.Now()).Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}
