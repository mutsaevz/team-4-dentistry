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
	logger *slog.Logger
}

func NewServService(
	services repository.ServiceRepository,
	logger *slog.Logger,
) ServService {
	return &servService{services: services,logger: logger}
}

func (s *servService) CreateService(
	req models.ServiceCreateRequest,
) (*models.Service, error) {
	if err := s.ValidateCreateServ(req); err != nil {
		return nil, err
	}

	service := &models.Service{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Category:    strings.TrimSpace(req.Category),
		Duration:    req.Duration,
		Price:       req.Price,
	}

	if err := s.services.Create(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *servService) GetServiceByID(id uint) (*models.Service, error) {
	service, err := s.services.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServiceNotfound
		}
		return nil, err
	}

	return service, nil
}

func (s *servService) ListServices(offset, limit int) ([]models.Service, error) {
	services, err := s.services.List(offset, limit)

	if err != nil {
		return nil, err
	}

	return services, err
}

func (s *servService) ListServicesByCategory(
	category string,
	offset,
	limit int) ([]models.Service, error) {
	services, err := s.services.ListByCategory(category, offset, limit)

	if err != nil {
		return nil, err
	}

	return services, nil
}

func (s *servService) UpdateService(
	id uint, req models.ServiceUpdateRequest,
) (*models.Service, error) {
	service, err := s.services.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServiceNotfound
		}
		return nil, err
	}

	if err := s.ApplyServUpdate(service, req); err != nil {
		return nil, err
	}

	if err := s.services.Update(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *servService) DeleteService(id uint) error {
	_, err := s.services.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrServiceNotfound
		}
		return err
	}

	if err := s.services.Delete(id); err != nil {
		return err
	}

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
