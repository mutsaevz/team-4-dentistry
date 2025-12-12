package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(*models.Review) error

	GetByID(context.Context, uint) (*models.Review, error)

	GetByDoctorID(context.Context, uint) ([]models.Review, error)

	GetByPatientID(context.Context, uint) ([]models.Review, error)

	Update(*models.Review) error

	Delete(context.Context, uint) error

	GetAverageRating(context.Context, uint) (float64, error)
}

type gormReviewRepository struct {
	DB *gorm.DB
	logger *slog.Logger
}

func NewReviewRepository(db *gorm.DB, logger *slog.Logger) ReviewRepository {
	return &gormReviewRepository{DB: db, logger: logger}
}

func (r *gormReviewRepository) Create(review *models.Review) error {
	if review == nil {
		return errors.New("review is nil")
	}

	return r.DB.Create(review).Error
}

func (r *gormReviewRepository) GetByID(ctx context.Context, id uint) (*models.Review, error) {
	var review models.Review
	if err := r.DB.WithContext(ctx).First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *gormReviewRepository) GetByDoctorID(ctx context.Context, doctorID uint) ([]models.Review, error) {
	var reviews []models.Review

	if err := r.DB.WithContext(ctx).Where("doctor_id = ?", doctorID).Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *gormReviewRepository) GetByPatientID(ctx context.Context, patient_id uint) ([]models.Review, error) {
	var reviews []models.Review

	if err := r.DB.WithContext(ctx).Where("patient_id = ?", patient_id).Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *gormReviewRepository) Update(req *models.Review) error {
	if req == nil {
		return nil
	}

	return r.DB.Save(req).Error
}

func (r *gormReviewRepository) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&models.Review{}, id).Error
}

func (r *gormReviewRepository) GetAverageRating(ctx context.Context, doctorID uint) (float64, error) {
	var avg sql.NullFloat64

	err := r.DB.WithContext(ctx).Model(&models.Review{}).
		Select("AVG(rating)").
		Where("doctor_id = ?", doctorID).
		Scan(&avg).Error
	if err != nil {
		return 0, err
	}

	if !avg.Valid {
		return 0, nil
	}

	return avg.Float64, nil
}
