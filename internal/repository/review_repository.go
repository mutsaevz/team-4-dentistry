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
	r.logger.Debug("создаем review в репозитории")
	if review == nil {
		r.logger.Error("передан nil review")
		return errors.New("review is nil")
	}

	if err := r.DB.Create(review).Error; err != nil {
		r.logger.Error("ошибка при создании review", "error", err)
		return err
	}

	r.logger.Info("review создан")
	return nil
}

func (r *gormReviewRepository) GetByID(ctx context.Context, id uint) (*models.Review, error) {
	r.logger.Debug("получаем review по ID в репозитории")
	var review models.Review
	if err := r.DB.WithContext(ctx).First(&review, id).Error; err != nil {
		r.logger.Error("ошибка при получении review по ID", "error", err)
		return nil, err
	}

	r.logger.Info("review получен по ID")
	return &review, nil
}

func (r *gormReviewRepository) GetByDoctorID(ctx context.Context, doctorID uint) ([]models.Review, error) {
	r.logger.Debug("получаем список review по doctorID в репозитории")
	var reviews []models.Review

	if err := r.DB.WithContext(ctx).Where("doctor_id = ?", doctorID).Find(&reviews).Error; err != nil {
		r.logger.Error("ошибка при получении списка review по doctorID", "error", err)
		return nil, err
	}

	r.logger.Info("список review получен по doctorID")

	return reviews, nil
}

func (r *gormReviewRepository) GetByPatientID(ctx context.Context, patient_id uint) ([]models.Review, error) {
	r.logger.Debug("получаем список review по patientID в репозитории")
	var reviews []models.Review

	if err := r.DB.WithContext(ctx).Where("id = ?", patient_id).Find(&reviews).Error; err != nil {
		r.logger.Error("ошибка при получении списка review по patientID", "error", err)
		return nil, err
	}

	r.logger.Info("список review получен по patientID")
	return reviews, nil
}

func (r *gormReviewRepository) Update(req *models.Review) error {
	r.logger.Debug("обновляем review в репозитории")
	if req == nil {
		r.logger.Error("передан nil review")
		return nil
	}

	r.logger.Info("review обновлен")

	return r.DB.Save(req).Error
}

func (r *gormReviewRepository) Delete(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&models.Review{}, id).Error; err != nil {
		r.logger.Error("ошибка при удалении review", "error", err)
		return err
	}

	r.logger.Info("review удален")
	return nil
}

func (r *gormReviewRepository) GetAverageRating(ctx context.Context, doctorID uint) (float64, error) {
	r.logger.Debug("получаем средний рейтинг по doctorID в репозитории")
	var avg sql.NullFloat64

	err := r.DB.WithContext(ctx).Model(&models.Review{}).
		Select("AVG(rating)").
		Where("doctor_id = ?", doctorID).
		Scan(&avg).Error
	if err != nil {
		r.logger.Error("ошибка при получении среднего рейтинга по doctorID", "error", err)
		return 0, err
	}

	if !avg.Valid {
		return 0, nil
	}

	r.logger.Info("средний рейтинг получен по doctorID")
	return avg.Float64, nil
}
