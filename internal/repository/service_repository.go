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
	db     *gorm.DB
	logger *slog.Logger
}

func NewServiceRepository(db *gorm.DB, logger *slog.Logger) ServiceRepository {
	return &gormServiceRepository{db: db, logger: logger}
}

func (r *gormServiceRepository) Create(service *models.Service) error {
	if service == nil {
		r.logger.Warn("попытка создать nil service")
		return constants.Service_IS_nil
	}

	r.logger.Debug("создание service в репозитории", "name", service.Name)

	if err := r.db.Create(service).Error; err != nil {
		r.logger.Error("ошибка при создании service", "error", err, "name", service.Name)
		return err
	}

	r.logger.Info("service создан", "service_id", service.ID, "name", service.Name)
	return nil
}

func (r *gormServiceRepository) GetByID(id uint) (*models.Service, error) {
	var service models.Service
	r.logger.Debug("получение service по ID", "service_id", id)

	if err := r.db.First(&service, id).Error; err != nil {
		r.logger.Error("ошибка при получении service по ID", "error", err, "service_id", id)
		return nil, err
	}

	r.logger.Info("service получен по ID", "service_id", id)
	return &service, nil
}

func (r *gormServiceRepository) List(offset, limit int) ([]models.Service, error) {
	var services []models.Service
	r.logger.Debug("получение списка services", "offset", offset, "limit", limit)

	if err := r.db.
		Offset(offset).
		Limit(limit).
		Find(&services).Error; err != nil {
		r.logger.Error("ошибка при получении списка services", "error", err, "offset", offset, "limit", limit)
		return nil, err
	}

	r.logger.Info("список services получен", "count", len(services))
	return services, nil
}

func (r *gormServiceRepository) ListByCategory(
	category string,
	offset,
	limit int,
) ([]models.Service, error) {
	var services []models.Service
	r.logger.Debug("получение services по категории", "category", category, "offset", offset, "limit", limit)

	if err := r.db.
		Where("category = ?", category).
		Offset(offset).
		Limit(limit).
		Find(&services).Error; err != nil {
		r.logger.Error("ошибка при получении services по категории", "error", err, "category", category)
		return nil, err
	}

	r.logger.Info("services по категории получены", "category", category, "count", len(services))
	return services, nil
}

func (r *gormServiceRepository) Update(service *models.Service) error {
	if service == nil {
		r.logger.Warn("попытка обновить nil service")
		return constants.Service_IS_nil
	}

	r.logger.Debug("обновление service", "service_id", service.ID)

	if err := r.db.Save(service).Error; err != nil {
		r.logger.Error("ошибка при обновлении service", "error", err, "service_id", service.ID)
		return err
	}

	r.logger.Info("service успешно обновлен", "service_id", service.ID)
	return nil
}

func (r *gormServiceRepository) Delete(id uint) error {
	r.logger.Debug("удаление service по ID", "service_id", id)

	if err := r.db.Delete(&models.Service{}, id).Error; err != nil {
		r.logger.Error("ошибка при удалении service", "error", err, "service_id", id)
		return err
	}

	r.logger.Info("service успешно удален", "service_id", id)
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
