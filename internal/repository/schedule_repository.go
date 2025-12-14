package repository

import (
	"context"
	"errors"
	"log/slog"
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
	DB     *gorm.DB
	logger *slog.Logger
}

func NewScheduleRepository(db *gorm.DB, logger *slog.Logger) ScheduleRepository {
	return &gormScheduleRepository{DB: db, logger: logger}
}

func (r *gormScheduleRepository) Create(ctx context.Context, schedule *models.Schedule) error {
	if schedule == nil {
		r.logger.Warn("попытка создать nil schedule")
		return errors.New("schedule is nil")
	}
	r.logger.Debug("создание schedule", "doctor_id", schedule.DoctorID, "date", schedule.Date)
	if err := r.DB.WithContext(ctx).Create(schedule).Error; err != nil {
		r.logger.Error("ошибка при создании schedule", "error", err)
		return err
	}
	r.logger.Info("schedule создан", "schedule_id", schedule.ID)
	return nil
}

func (r *gormScheduleRepository) GetAll(ctx context.Context) ([]models.Schedule, error) {
	r.logger.Debug("получение всех schedules")
	var schedules []models.Schedule

	if err := r.DB.WithContext(ctx).Find(&schedules).Error; err != nil {
		r.logger.Error("ошибка при получении всех schedules", "error", err)
		return nil, err
	}

	r.logger.Info("успешное получение всех schedules", "count", len(schedules))
	return schedules, nil
}

func (r *gormScheduleRepository) GetByID(ctx context.Context, id uint) (*models.Schedule, error) {
	r.logger.Debug("получение schedule по ID", "schedule_id", id)
	var schedule models.Schedule

	if err := r.DB.WithContext(ctx).First(&schedule, id).Error; err != nil {
		r.logger.Error("ошибка при получении schedule по ID", "error", err, "schedule_id", id)
		return nil, err
	}

	r.logger.Info("schedule получен по ID", "schedule_id", id)
	return &schedule, nil
}

func (r *gormScheduleRepository) GetSchedulesByDoctorID(ctx context.Context, doctorID uint) ([]models.Schedule, error) {
	r.logger.Debug("получение schedules по doctorID", "doctor_id", doctorID)
	var schedule []models.Schedule

	if err := r.DB.WithContext(ctx).Where("doctor_id = ?", doctorID).Find(&schedule).Error; err != nil {
		r.logger.Error("ошибка при получении schedules по doctorID", "error", err, "doctor_id", doctorID)
		return nil, err
	}

	r.logger.Info("schedules получены по doctorID", "doctor_id", doctorID, "count", len(schedule))
	return schedule, nil
}

func (r *gormScheduleRepository) Update(ctx context.Context, schedule *models.Schedule) error {
	if schedule == nil {
		r.logger.Warn("попытка обновить nil schedule")
		return nil
	}
	r.logger.Debug("обновление schedule", "schedule_id", schedule.ID)
	if err := r.DB.WithContext(ctx).Save(schedule).Error; err != nil {
		r.logger.Error("ошибка при обновлении schedule", "error", err, "schedule_id", schedule.ID)
		return err
	}

	r.logger.Info("schedule успешно обновлен", "schedule_id", schedule.ID)
	return nil
}

func (r *gormScheduleRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Debug("удаление schedule по ID", "schedule_id", id)
	if err := r.DB.WithContext(ctx).Delete(&models.Schedule{}, id).Error; err != nil {
		r.logger.Error("ошибка при удалении schedule", "error", err, "schedule_id", id)
		return err
	}
	r.logger.Info("schedule успешно удален", "schedule_id", id)
	return nil
}

func (r *gormScheduleRepository) DeleteByDoctorID(ctx context.Context, doctorID uint) error {
	r.logger.Debug("удаление schedules по doctorID", "doctor_id", doctorID)
	if err := r.DB.WithContext(ctx).Where("doctor_id = ?", doctorID).Delete(&models.Schedule{}).Error; err != nil {
		r.logger.Error("ошибка при удалении schedules по doctorID", "error", err, "doctor_id", doctorID)
		return err
	}
	r.logger.Info("schedules успешно удалены по doctorID", "doctor_id", doctorID)
	return nil
}

func (r *gormScheduleRepository) GetAvailableSlots(
	ctx context.Context,
	doctorID uint,
	start, end time.Time,
) ([]models.Schedule, error) {

	start = start.Truncate(24 * time.Hour)
	end = start.AddDate(0, 0, 7)

	var schedules []models.Schedule

	r.logger.Debug("получение доступных слотов", "doctor_id", doctorID, "start", start, "end", end)

	err := r.DB.WithContext(ctx).
		Where("doctor_id = ?", doctorID).
		Where("is_available = ?", true).
		Where("start_time >= ? AND end_time <= ?", start, end).
		Order("start_time ASC").
		Find(&schedules).Error

	if err != nil {
		r.logger.Error("ошибка при получении доступных слотов", "error", err)
		return nil, err
	}

	r.logger.Info("доступные слоты получены", "doctor_id", doctorID, "count", len(schedules))
	return schedules, nil
}
