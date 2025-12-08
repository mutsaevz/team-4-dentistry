package repository

import (
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(service *models.Service) error

	GetByID(id uint) (*models.Service, error)

	Update(service *models.Service) error

	Delete(id uint) error
}

type gormServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &gormServiceRepository{db: db}
}

func (r *gormServiceRepository) Create(service *models.Service) error {
	if service == nil {
		return nil
	}

	if err := r.db.Create(&service).Error; err != nil {
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

func (r *gormServiceRepository) Update(service *models.Service) error {
	if service == nil {
		return nil
	}

	if err := r.db.Save(&service).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormServiceRepository) Delete(id uint) error {
	var service models.Service

	if err := r.db.Delete(&service, id).Error; err != nil {
		return err
	}

	return nil
}
