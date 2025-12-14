package services

import (
	"errors"
	"log/slog"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrRecommendationNotFound = errors.New("рекомендация не найдена")
	ErrInvalidPatiendID       = errors.New("некорректный patient_id")
	ErrInvalidServiceId       = errors.New("некорректный service_id")
)

type RecommendationService interface {
	CreateRec(doctorID uint, req models.RecommendationCreateRequest) (*models.Recommendation, error)

	ListRecsByPatientID(patientID uint) ([]models.Recommendation, error)

	DeleteRec(id uint) error
}

type recommendationService struct {
	recRepo     repository.RecommendationRepository
	userRepo    repository.UserRepository
	serviceRepo repository.ServiceRepository
	logger      *slog.Logger
}

func NewRecommendationService(
	recRepo repository.RecommendationRepository,
	userRepo repository.UserRepository,
	serviceRepo repository.ServiceRepository,
	logger *slog.Logger,
) *recommendationService {
	return &recommendationService{
		recRepo:     recRepo,
		userRepo:    userRepo,
		serviceRepo: serviceRepo,
		logger:      logger,
	}
}

func (s *recommendationService) CreateRec(
	doctorID uint,
	req models.RecommendationCreateRequest,
) (*models.Recommendation, error) {
	s.logger.Debug("CreateRec вызван", "doctor_id", doctorID, "patient_id", req.PatientID, "service_id", req.ServiceID)
	if err := s.ValidateCreate(req); err != nil {
		s.logger.Error("валидация CreateRec провалилась", "error", err)
		return nil, err
	}

	if _, err := s.userRepo.GetByID(req.PatientID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidPatiendID
		}
		return nil, err
	}

	service, err := s.serviceRepo.GetByID(req.ServiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidServiceId
		}
		return nil, err
	}

	rec := &models.Recommendation{
		PatientID: req.PatientID,
		ServiceID: req.ServiceID,
		DoctorID:  doctorID,
		Note:      req.Note,
		Service:   service,
	}

	if err := s.recRepo.Create(rec); err != nil {
		s.logger.Error("ошибка при создании recommendation", "error", err)
		return nil, err
	}

	s.logger.Info("recommendation создан", "id", rec.ID, "doctor_id", doctorID, "patient_id", rec.PatientID)
	return rec, nil
}

func (s *recommendationService) ListRecsByPatientID(
	patientID uint,
) ([]models.Recommendation, error) {
	if patientID <= 0 {
		return nil, ErrInvalidPatiendID
	}
	s.logger.Debug("ListRecsByPatientID вызван", "patient_id", patientID)

	recs, err := s.recRepo.ListByPatientID(patientID)
	if err != nil {
		s.logger.Error("ошибка при получении рекомендаций по patient_id", "error", err, "patient_id", patientID)
		return nil, err
	}

	s.logger.Info("список рекомендаций получен", "patient_id", patientID, "count", len(recs))
	return recs, nil
}

func (s *recommendationService) DeleteRec(id uint) error {
	if id <= 0 {
		return ErrRecommendationNotFound
	}
	s.logger.Debug("DeleteRec вызван", "id", id)

	_, err := s.recRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("recommendation не найдена при DeleteRec", "id", id)
			return ErrRecommendationNotFound
		}
		s.logger.Error("ошибка при получении recommendation перед удалением", "error", err, "id", id)
		return err
	}

	if err := s.recRepo.Delete(id); err != nil {
		s.logger.Error("ошибка при удалении recommendation", "error", err, "id", id)
		return err
	}

	s.logger.Info("recommendation удален", "id", id)
	return nil
}

func (s *recommendationService) ValidateCreate(
	req models.RecommendationCreateRequest,
) error {
	if req.PatientID <= 0 {
		return ErrInvalidPatiendID
	}

	if req.ServiceID <= 0 {
		return ErrInvalidServiceId
	}

	s.logger.Debug("ValidateCreate успешно", "patient_id", req.PatientID, "service_id", req.ServiceID)
	return nil
}
