package repository

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

var validate = validator.New()

type DoctorRepository interface {
	Create(*models.Doctor) error

	GetAll(models.DoctorQueryParams, context.Context) ([]models.Doctor, error)

	GetByID(uint, context.Context) (*models.Doctor, error)

	Update(*models.Doctor) error

	UpdateAvgRating(uint, float64) error

	Delete(uint) error
}

type gormDoctorRepository struct {
	DB *gorm.DB
}

func NewDoctorRepositry(db *gorm.DB) DoctorRepository {
	return &gormDoctorRepository{DB: db}
}

func (r *gormDoctorRepository) Create(doctor *models.Doctor) error {
	if doctor == nil {
		return errors.New("doctor is nil")
	}

	return r.DB.Create(doctor).Error
}

func (r *gormDoctorRepository) GetAll(params models.DoctorQueryParams, ctx context.Context) ([]models.Doctor, error) {
	var doctors []models.Doctor

	q := r.DB.Model(&models.Doctor{})

	if err := validate.Struct(params); err != nil {
		return nil, err
	}

	if params.FilOr {
		q = q.Where("specialization ILIKE ? OR experience_years >= ? OR avg_rating >= ?",
			"%"+params.Specialization+"%",
			params.ExperienceYears,
			params.AvgRating)
	} else {
		if params.Specialization != "" {
			q = r.DB.Where("specialization ILIKE ?", "%"+params.Specialization+"%")
		}

		if params.ExperienceYears > 0 {
			q = r.DB.Where("experience_years >= ?", params.ExperienceYears)
		}

		if params.AvgRating > 0 {
			q = r.DB.Where("avg_rating >= ?", params.AvgRating)
		}
	}

	if err := q.WithContext(ctx).Find(&doctors).Error; err != nil {
		return nil, err
	}

	return doctors, nil
}

func (r *gormDoctorRepository) GetByID(id uint, ctx context.Context) (*models.Doctor, error) {
	var doctor models.Doctor

	if err := r.DB.WithContext(ctx).First(&doctor, id).Error; err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (r *gormDoctorRepository) Update(doctor *models.Doctor) error {
	if doctor == nil {
		return nil
	}

	return r.DB.Save(doctor).Error
}

func (r *gormDoctorRepository) UpdateAvgRating(id uint, avg float64) error {
	return r.DB.Model(&models.Doctor{}).
		Where("id = ?", id).
		Update("avg_rating", avg).Error
}

func (r *gormDoctorRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Doctor{}, id).Error
}
