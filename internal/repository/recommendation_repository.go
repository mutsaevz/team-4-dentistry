package repository

import (
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
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) RecommendationRepository {
	return &gormRecommendationRepository{db: db}
}

func (r *gormRecommendationRepository) Create(rec *models.Recommendation) error {
	if rec == nil {
		return constants.Rec_IS_nil
	}

	if err := r.db.Create(rec).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormRecommendationRepository) GetByID(id uint) (*models.Recommendation, error) {
	var rec models.Recommendation

	if err := r.db.
		Preload("Service").
		Preload("Patient").
		First(&rec, id).Error; err != nil {
		return nil, err
	}

	return &rec, nil
}

func (r *gormRecommendationRepository) ListByPatientID(
	patientID uint,
) ([]models.Recommendation, error) {
	var recs []models.Recommendation

	if err := r.db.
		Preload("Service").
		Preload("Patient").
		Where("patient_id = ?", patientID).
		Order("created_at DESC").
		Find(&recs).Error; err != nil {
		return nil, err
	}

	return recs, nil
}

func (r *gormRecommendationRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Recommendation{}, id).Error; err != nil {
		return err
	}

	return nil
}
