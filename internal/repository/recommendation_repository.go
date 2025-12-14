package repository

import (
	"log/slog"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type RecommendationRepository interface {
	Create(rec *models.Recommendation) error

	GetByID(id uint) (*models.Recommendation, error)

	ListByPatientID(patientID uint) ([]models.Recommendation, error)

	Delete(id uint) error
}

type gormRecommendationRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRecommendationRepository(db *gorm.DB, logger *slog.Logger) RecommendationRepository {
	return &gormRecommendationRepository{db: db, logger: logger}
}

func (r *gormRecommendationRepository) Create(rec *models.Recommendation) error {
	r.logger.Debug("cоздаем recommendation в репозитории")

	if rec == nil {

		r.logger.Error("передан nil recommendation")
		return constants.Rec_IS_nil
	}

	if err := r.db.Create(rec).Error; err != nil {
		r.logger.Error("ошибка при создании recommendation", "error", err)
		return err
	}

	r.logger.Info("recommendation создан")

	return nil
}

func (r *gormRecommendationRepository) GetByID(id uint) (*models.Recommendation, error) {

	r.logger.Debug("получаем recommendation по ID в репозитории")
	var rec models.Recommendation

	if err := r.db.
		Preload("Service").
		Preload("Patient").
		First(&rec, id).Error; err != nil {
		r.logger.Error("ошибка при получении recommendation по ID", "error", err)
		return nil, err
	}

	r.logger.Info("recommendation получен по ID")

	return &rec, nil
}

func (r *gormRecommendationRepository) ListByPatientID(
	patientID uint,
) ([]models.Recommendation, error) {

	r.logger.Debug("получаем список recommendation по patientID в репозитории")
	var recs []models.Recommendation

	if err := r.db.
		Preload("Service").
		Preload("Patient").
		Where("patient_id = ?", patientID).
		Order("created_at DESC").
		Find(&recs).Error; err != nil {

		r.logger.Error("ошибка при получении списка recommendation по patientID", "error", err)

		return nil, err
	}

	r.logger.Info("список recommendation получен по patientID")

	return recs, nil
}

func (r *gormRecommendationRepository) Delete(id uint) error {

	r.logger.Debug("удаляем recommendation в репозитории")
	if err := r.db.Delete(&models.Recommendation{}, id).Error; err != nil {
		r.logger.Error("ошибка при удалении recommendation", "error", err)
		return err
	}

	r.logger.Info("recommendation удален")
	return nil
}
