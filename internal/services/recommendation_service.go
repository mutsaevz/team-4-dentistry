package services

import (
	"errors"

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
}

func NewRecommendationService(
	recRepo repository.RecommendationRepository,
	userRepo repository.UserRepository,
	serviceRepo repository.ServiceRepository,
) *recommendationService {
	return &recommendationService{
		recRepo:     recRepo,
		userRepo:    userRepo,
		serviceRepo: serviceRepo,
	}
}

func (s *recommendationService) CreateRec(
	doctorID uint,
	req models.RecommendationCreateRequest,
) (*models.Recommendation, error) {
	if err := s.ValidateCreate(req); err != nil {
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
		return nil, err
	}

	return rec, nil
}

func (s *recommendationService) ListRecsByPatientID(
	patientID uint,
) ([]models.Recommendation, error) {

	if patientID <= 0 {
		return nil, ErrInvalidPatiendID
	}

	recs, err := s.recRepo.ListByPatientID(patientID)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (s *recommendationService) DeleteRec(id uint) error {
	if id <= 0 {
		return ErrRecommendationNotFound
	}

	_, err := s.recRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecommendationNotFound
		}
		return err
	}

	if err := s.recRepo.Delete(id); err != nil {
		return err
	}
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

	return nil
}
