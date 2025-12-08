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

	GetAvailableSlots(uint) ([]models.Shedule, error)
}

type gormSheduleRepository struct {
	DB *gorm.DB
}

func NewSheduleRepository(db *gorm.DB) SheduleRepositroy {
	return &gormSheduleRepository{DB: db}
}

func (r *gormSheduleRepository) Create(shedule *models.Shedule) error {
	if shedule == nil {
		return errors.New("shedule is nil")
	}

	return r.DB.Create(shedule).Error
}

func (r *gormSheduleRepository) GetAll(ctx context.Context) ([]models.Shedule, error) {
	var shedules []models.Shedule

	if err := r.DB.WithContext(ctx).Find(&shedules).Error; err != nil {
		return nil, err
	}

	return shedules, nil
}

func (r *gormSheduleRepository) GetByID(id uint, ctx context.Context) (*models.Shedule, error) {
	var shedule models.Shedule

	if err := r.DB.WithContext(ctx).First(&shedule).Error; err != nil {
		return nil, err
	}

	return &shedule, nil
}

func (r *gormSheduleRepository) GetByDateRange(ctx context.Context, doctorID uint, start, end time.Time) (*models.Shedule, error) {
	var shedule models.Shedule

	if err := r.DB.WithContext(ctx).Where("start_time <= ? AND end_time >= ?", start, end).First(&shedule, doctorID).Error; err != nil {
		return nil, err
	}

	return &shedule, nil
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

func (r *gormSheduleRepository) Update(shedule *models.Shedule) error {
	if shedule == nil {
		return nil
	}

	return r.DB.Save(shedule).Error
}

func (r *gormSheduleRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Shedule{}, id).Error
}

func (r *gormSheduleRepository) DeleteByDoctorID(doctorID uint) error {
	return r.DB.Where("doctor_id = ?", doctorID).Delete(&models.Shedule{}).Error
}

func (r *gormSheduleRepository) GetAvailableSlots(doctorID uint) ([]models.Shedule, error) {
	var slots []models.Shedule

	subQuery := r.DB.Model(&models.Shedule{}).
		Select("schedule_id").
		Where("status != ?", "canceled")

	if err := r.DB.Where("doctor_id = ? AND id NOT IN (?)", doctorID, subQuery).
		Order("start_time").
		Find(&slots).Error; err != nil {
		return nil, err
	}

	return slots, nil
}
