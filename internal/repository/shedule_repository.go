package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type SheduleRepositroy interface {
	Create(*models.Shedule) error

	GetAll(context.Context) ([]models.Shedule, error)

	GetByID(uint, context.Context) (*models.Shedule, error)

	GetByDateRange(context.Context, uint, time.Time, time.Time) (*models.Shedule, error)

	CheckConflict(uint, time.Time, time.Time) (bool, error)

	Update(*models.Shedule) error

	Delete(uint) error

	DeleteByDoctorID(uint) error

	GetAvailableSlots(context.Context, uint) ([]models.Shedule, error)
}

type gormSheduleRepository struct {
	DB *gorm.DB
}

func NewSheduleRepository(db *gorm.DB) SheduleRepositroy {
	return &gormSheduleRepository{DB: db}
}

func (r *gormSheduleRepository) Create(schedule *models.Shedule) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}

	return r.DB.Create(schedule).Error
}

func (r *gormSheduleRepository) GetAll(ctx context.Context) ([]models.Shedule, error) {
	var schedules []models.Shedule

	if err := r.DB.WithContext(ctx).Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *gormSheduleRepository) GetByID(id uint, ctx context.Context) (*models.Shedule, error) {
	var schedule models.Shedule

	if err := r.DB.WithContext(ctx).First(&schedule, id).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *gormSheduleRepository) GetByDateRange(ctx context.Context, doctorID uint, start, end time.Time) (*models.Shedule, error) {
	var schedule models.Shedule

	if err := r.DB.WithContext(ctx).Where("start_time <= ? AND end_time >= ?", start, end).First(&schedule, doctorID).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *gormSheduleRepository) CheckConflict(doctorID uint, startTime, endTime time.Time) (bool, error) {
	var count int64

	err := r.DB.Model(&models.Shedule{}).
		Where("doctor_id = ? AND start_time < ? AND end_time > ?", doctorID, endTime, startTime).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *gormSheduleRepository) Update(schedule *models.Shedule) error {
	if schedule == nil {
		return nil
	}

	return r.DB.Save(schedule).Error
}

func (r *gormSheduleRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Shedule{}, id).Error
}

func (r *gormSheduleRepository) DeleteByDoctorID(doctorID uint) error {
	return r.DB.Where("doctor_id = ?", doctorID).Delete(&models.Shedule{}).Error
}

func (r *gormSheduleRepository) GetAvailableSlots(ctx context.Context, doctorID uint) ([]models.Shedule, error) {
	var slots []models.Shedule

	subQuery := r.DB.Model(&models.Appointment{}).
		Select("schedule_id").
		Where("status != ?", "canceled")

	if err := r.DB.WithContext(ctx).Where("doctor_id = ? AND id NOT IN (?)", doctorID, subQuery).
		Order("start_time").
		Find(&slots).Error; err != nil {
		return nil, err
	}

	return slots, nil
}
