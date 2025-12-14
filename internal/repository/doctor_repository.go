package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

var validate = validator.New()

type DoctorRepository interface {
	Create(context.Context, *models.Doctor) error

	GetAll(models.DoctorQueryParams, context.Context) ([]models.Doctor, error)

	GetByID(uint, context.Context) (*models.Doctor, error)

	Update(context.Context, *models.Doctor) error

	UpdateAvgRating(context.Context, uint, float64) error

	Delete(context.Context, uint) error
}

type gormDoctorRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewDoctorRepository(db *gorm.DB, logger *slog.Logger) DoctorRepository {
	return &gormDoctorRepository{DB: db, logger: logger}
}

func (r *gormDoctorRepository) Create(ctx context.Context, doctor *models.Doctor) error {
	if doctor == nil {
		r.logger.Warn("doctor равен nil")
		return errors.New("doctor is nil")
	}

	if err := r.DB.WithContext(ctx).Create(doctor).Error; err != nil {
		r.logger.Error("ошибка при создании doctor", "ошибка", err)
	}

	r.logger.Info("doctor успешно создан", "doctor_id", doctor.ID)
	return nil
}

func (r *gormDoctorRepository) GetAll(params models.DoctorQueryParams, ctx context.Context) ([]models.Doctor, error) {
	r.logger.Debug("Получение всех doctors с параметрами", "params", params)
	var doctors []models.Doctor

	q := r.DB.WithContext(ctx).Model(&models.Doctor{})

	if err := validate.Struct(params); err != nil {
		r.logger.Error("ошибка валидации параметров запроса doctors", "ошибка", err, "params", params)
		return nil, err
	}

	if params.FilOr {
		q = q.Where("specializations ILIKE ? OR experience_years >= ? OR avg_rating >= ?",
			"%"+params.Specialization+"%",
			params.ExperienceYears,
			params.AvgRating)
	} else {
		if params.Specialization != "" {
			q = q.Where("specializations ILIKE ?", "%"+params.Specialization+"%")
		}

		if params.ExperienceYears > 0 {
			q = q.Where("experience_years >= ?", params.ExperienceYears)
		}

		if params.AvgRating > 0 {
			q = q.Where("avg_rating >= ?", params.AvgRating)
		}
	}

	if err := q.Find(&doctors).Error; err != nil {
		r.logger.Error("ошибка при получении всех doctors", "ошибка", err, "params", params)
		return nil, err
	}

	r.logger.Info("успешное получение всех doctors", "params", params)

	return doctors, nil
}

func (r *gormDoctorRepository) GetByID(id uint, ctx context.Context) (*models.Doctor, error) {

	r.logger.Debug("Получение doctor по ID", "doctor_id", id)
	var doctor models.Doctor

	if err := r.DB.WithContext(ctx).First(&doctor, id).Error; err != nil {
		r.logger.Error("ошибка при получении doctor по ID", "ошибка", err, "doctor_id", id)
		return nil, err
	}

	r.logger.Info("успешное получение doctor по ID", "doctor_id", id)

	return &doctor, nil
}

func (r *gormDoctorRepository) Update(ctx context.Context, doctor *models.Doctor) error {
	if doctor == nil {
		r.logger.Warn("doctor равен nil")
		return errors.New("")
	}

	if err := r.DB.WithContext(ctx).Save(doctor).Error; err != nil {
		r.logger.Error("ошибка при обновлении doctor", "ошибка", err, "doctor_id", doctor.ID)
		return err
	}

	r.logger.Info("успешное обновление doctor", "doctor_id", doctor.ID)
	return nil
}

func (r *gormDoctorRepository) UpdateAvgRating(ctx context.Context, id uint, avg float64) error {

	r.logger.Debug("Обновление среднего рейтинга doctor", "doctor_id", id, "avg_rating", avg)

	if err := r.DB.WithContext(ctx).Model(&models.Doctor{}).
		Where("id = ?", id).
		Update("avg_rating", avg).Error; err != nil {
		r.logger.Error("ошибка при обновлении doctor", "error", err)
		return err
	}

	r.logger.Info("успешное обновление среднего рейтинга doctor", "doctor_id", id, "avg_rating", avg)
	return nil
}

func (r *gormDoctorRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Debug("Удаление doctor по ID", "doctor_id", id)
	if err := r.DB.WithContext(ctx).Delete(&models.Doctor{}, id).Error; err != nil {
		r.logger.Error("ошибка при удалении doctor", "error", err, "doctor_id", id)
		return err
	}
	r.logger.Info("успешное удаление doctor", "doctor_id", id)
	return nil
}
