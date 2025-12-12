package repository

import (
	"context"
	"log/slog"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(service *models.Service) error

	GetByID(id uint) (*models.Service, error)

	List(offset, limit int) ([]models.Service, error)

	ListByCategory(category string, offset, limit int) ([]models.Service, error)

	Update(service *models.Service) error

	Delete(id uint) error

	GetServicesByDoctorID(ctx context.Context, doctorID uint) ([]models.Service, error)
}

type gormServiceRepository struct {
	db *gorm.DB
	logger *slog.Logger
}

func NewServiceRepository(db *gorm.DB, logger *slog.Logger) ServiceRepository {
	return &gormServiceRepository{db: db, logger: logger}
}

func (r *gormServiceRepository) Create(service *models.Service) error {
	if service == nil {
		return constants.Service_IS_nil
	}

	if err := r.db.Create(service).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormServiceRepository) GetByID(id uint) (*models.Service, error) {
	var service models.Service

	if err := r.db.First(&service, id).Error; err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *gormServiceRepository) List(offset, limit int) ([]models.Service, error) {
	var services []models.Service

	if err := r.db.
		Offset(offset).
		Limit(limit).
		Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}

func (r *gormServiceRepository) ListByCategory(
	category string,
	offset,
	limit int,
) ([]models.Service, error) {
	var services []models.Service

	if err := r.db.
		Where("category = ?", category).
		Offset(offset).
		Limit(limit).
		Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}

func (r *gormServiceRepository) Update(service *models.Service) error {
	if service == nil {
		return constants.Service_IS_nil
	}

	if err := r.db.Save(service).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormServiceRepository) Delete(id uint) error {

	if err := r.db.Delete(&models.Service{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormServiceRepository) GetServicesByDoctorID(ctx context.Context, doctorID uint) ([]models.Service, error) {
	var services []models.Service

	if err := r.db.
		WithContext(ctx).
		Where("doctor_id = ?", doctorID).
		Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}
