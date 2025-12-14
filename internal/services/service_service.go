package services

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
	"gorm.io/gorm"
)

var ErrServiceNotfound = errors.New("услуга не найдена")

type ServService interface {
	CreateService(req models.ServiceCreateRequest) (*models.Service, error)

	GetServiceByID(id uint) (*models.Service, error)

	ListServices(offset, limit int) ([]models.Service, error)

	ListServicesByCategory(category string, offset, limit int) ([]models.Service, error)

	UpdateService(id uint, req models.ServiceUpdateRequest) (*models.Service, error)

	DeleteService(id uint) error
}

type servService struct {
	services repository.ServiceRepository
	logger   *slog.Logger
}

func NewServService(
	services repository.ServiceRepository,
	logger *slog.Logger,
) ServService {
	return &servService{services: services, logger: logger}
}

func (s *servService) CreateService(
	req models.ServiceCreateRequest,
) (*models.Service, error) {
	s.logger.Debug("CreateService called", "name", req.Name, "doctor_id", req.DoctorID)

	if err := s.ValidateCreateServ(req); err != nil {
		s.logger.Error("validation failed for CreateService", "error", err)
		return nil, err
	}

	service := &models.Service{
		DoctorID:    req.DoctorID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Category:    strings.TrimSpace(req.Category),
		Duration:    req.Duration,
		Price:       req.Price,
	}

	if err := s.services.Create(service); err != nil {
		s.logger.Error("failed to create service in repo", "error", err, "name", req.Name)
		return nil, err
	}

	s.logger.Info("service created", "service_id", service.ID, "name", service.Name)
	return service, nil
}

func (s *servService) GetServiceByID(id uint) (*models.Service, error) {
	s.logger.Debug("GetServiceByID called", "service_id", id)
	service, err := s.services.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("service not found", "service_id", id)
			return nil, ErrServiceNotfound
		}
		s.logger.Error("error getting service by id", "error", err, "service_id", id)
		return nil, err
	}

	s.logger.Info("service retrieved", "service_id", id)
	return service, nil
}

func (s *servService) ListServices(offset, limit int) ([]models.Service, error) {
	s.logger.Debug("ListServices called", "offset", offset, "limit", limit)
	services, err := s.services.List(offset, limit)
	if err != nil {
		s.logger.Error("error listing services", "error", err)
		return nil, err
	}
	s.logger.Info("services listed", "count", len(services))
	return services, nil
}

func (s *servService) ListServicesByCategory(
	category string,
	offset,
	limit int) ([]models.Service, error) {
	s.logger.Debug("ListServicesByCategory called", "category", category, "offset", offset, "limit", limit)
	services, err := s.services.ListByCategory(category, offset, limit)
	if err != nil {
		s.logger.Error("error listing services by category", "error", err, "category", category)
		return nil, err
	}
	s.logger.Info("services by category listed", "category", category, "count", len(services))
	return services, nil
}

func (s *servService) UpdateService(
	id uint, req models.ServiceUpdateRequest,
) (*models.Service, error) {
	s.logger.Debug("UpdateService called", "service_id", id)
	service, err := s.services.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("service not found for update", "service_id", id)
			return nil, ErrServiceNotfound
		}
		s.logger.Error("error fetching service for update", "error", err, "service_id", id)
		return nil, err
	}

	if err := s.ApplyServUpdate(service, req); err != nil {
		s.logger.Error("validation failed when applying service update", "error", err, "service_id", id)
		return nil, err
	}

	if err := s.services.Update(service); err != nil {
		s.logger.Error("failed to update service in repo", "error", err, "service_id", id)
		return nil, err
	}

	s.logger.Info("service updated", "service_id", id)
	return service, nil
}

func (s *servService) DeleteService(id uint) error {
	s.logger.Debug("DeleteService called", "service_id", id)
	_, err := s.services.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("service not found for delete", "service_id", id)
			return ErrServiceNotfound
		}
		s.logger.Error("error fetching service for delete", "error", err, "service_id", id)
		return err
	}

	if err := s.services.Delete(id); err != nil {
		s.logger.Error("failed to delete service in repo", "error", err, "service_id", id)
		return err
	}

	s.logger.Info("service deleted", "service_id", id)
	return nil
}

func (s *servService) ValidateCreateServ(req models.ServiceCreateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("название не должно быть пустым")
	}

	if strings.TrimSpace(req.Category) == "" {
		return errors.New("категория не должна быть пустой")
	}

	if req.Duration < 0 {
		return errors.New("время не должно быть отрицательным")
	}

	if req.Price < 0 {
		return errors.New("цена не должна быть отрицательной")
	}

	if req.DoctorID == 0 {
		return errors.New("doctor_id обязателен")
	}

	return nil
}

func (s *servService) ApplyServUpdate(
	service *models.Service, req models.ServiceUpdateRequest,
) error {
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)

		if trimmed == "" {
			return errors.New("название должно быть обязательно")
		}
		service.Name = trimmed
	}

	if req.Price != nil {
		if *req.Price < 0 {
			return errors.New("цена не должна быть отрицательной")
		}
		service.Price = *req.Price
	}

	if req.Duration != nil {
		if *req.Duration < 0 {
			return errors.New("время не должно быть отрицательной")
		}

		service.Duration = *req.Duration
	}

	if req.Category != nil {
		trimmed := strings.TrimSpace(*req.Category)
		if trimmed == "" {
			return errors.New("категория не должна быть пустой")
		}
		service.Category = trimmed
	}
	return nil
}
