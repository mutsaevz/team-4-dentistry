package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
)

type ReviewService interface {
	CreateReview(context.Context, models.ReviewCreateRequest) (*models.Review, error)

	GetByID(context.Context, uint) (*models.Review, error)

	UpdateReview(context.Context, uint, models.ReviewUpdateRequest) (*models.Review, error)

	DeleteReview(context.Context, uint) error

	GetDoctorReviews(context.Context, uint) ([]models.Review, error)

	GetPatientReviews(context.Context, uint) ([]models.Review, error)
}

type reviewService struct {
	review  repository.ReviewRepository
	doctor  repository.DoctorRepository
	patient repository.UserRepository
	logger  *slog.Logger
}

func NewReviewService(review repository.ReviewRepository,
	doctor repository.DoctorRepository,
	patient repository.UserRepository,
	logger *slog.Logger) ReviewService {
	return &reviewService{
		review:  review,
		doctor:  doctor,
		patient: patient,
		logger:  logger,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, req models.ReviewCreateRequest) (*models.Review, error) {
	s.logger.Debug("CreateReview вызван", "appointment_id", req.AppointmentID, "user_id", req.UserID, "doctor_id", req.DoctorID)

	if err := s.ValidateCreateReview(req); err != nil {
		s.logger.Error("валидация CreateReview провалилась", "error", err)
		return nil, err
	}

	var review = models.Review{
		AppointmentID: req.AppointmentID,
		UserID:        req.UserID,
		DoctorID:      req.DoctorID,
		Rating:        req.Rating,
		Comment:       req.Comment,
	}

	if err := s.review.Create(&review); err != nil {
		s.logger.Error("ошибка при создании review", "error", err)
		return nil, err
	}

	s.logger.Info("review создан", "review_id", review.ID, "appointment_id", review.AppointmentID)
	return &review, nil
}

func (s *reviewService) GetByID(ctx context.Context, id uint) (*models.Review, error) {
	s.logger.Debug("GetByID review вызван", "review_id", id)
	rev, err := s.review.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("ошибка получения review по ID", "error", err, "review_id", id)
		return nil, err
	}
	s.logger.Info("review получен по ID", "review_id", id)
	return rev, nil
}

func (s *reviewService) UpdateReview(ctx context.Context, id uint, req models.ReviewUpdateRequest) (*models.Review, error) {
	s.logger.Debug("UpdateReview вызван", "review_id", id)
	review, err := s.review.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("ошибка получения review для обновления", "error", err, "review_id", id)
		return nil, err
	}

	if req.Rating != nil {
		review.Rating = *req.Rating
	}

	if req.Comment != nil {
		review.Comment = *req.Comment
	}

	if err := s.review.Update(review); err != nil {
		s.logger.Error("ошибка при обновлении review", "error", err, "review_id", id)
		return nil, err
	}

	s.logger.Info("review успешно обновлен", "review_id", id)
	return review, nil
}

func (s *reviewService) DeleteReview(ctx context.Context, id uint) error {
	s.logger.Debug("DeleteReview вызван", "review_id", id)
	if err := s.review.Delete(ctx, id); err != nil {
		s.logger.Error("ошибка при удалении review", "error", err, "review_id", id)
		return err
	}
	s.logger.Info("review удален", "review_id", id)
	return nil
}

func (s *reviewService) GetDoctorReviews(ctx context.Context, doctorID uint) ([]models.Review, error) {
	s.logger.Debug("GetDoctorReviews вызван", "doctor_id", doctorID)
	revs, err := s.review.GetByDoctorID(ctx, doctorID)
	if err != nil {
		s.logger.Error("ошибка при получении отзывов врача", "error", err, "doctor_id", doctorID)
		return nil, err
	}
	s.logger.Info("отзывы врача получены", "doctor_id", doctorID, "count", len(revs))
	return revs, nil
}

func (s *reviewService) GetPatientReviews(ctx context.Context, patientID uint) ([]models.Review, error) {
	s.logger.Debug("GetPatientReviews вызван", "patient_id", patientID)
	revs, err := s.review.GetByPatientID(ctx, patientID)
	if err != nil {
		s.logger.Error("ошибка при получении отзывов пациента", "error", err, "patient_id", patientID)
		return nil, err
	}
	s.logger.Info("отзывы пациента получены", "patient_id", patientID, "count", len(revs))
	return revs, nil
}

func (s *reviewService) ValidateCreateReview(req models.ReviewCreateRequest) error {
	if req.AppointmentID == 0 {
		return errors.New("")
	}
	if req.UserID == 0 {
		return errors.New("user_id is required")
	}
	if req.DoctorID == 0 {
		return errors.New("doctor_id is required")
	}
	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	s.logger.Debug("ValidateCreateReview успешно", "appointment_id", req.AppointmentID, "user_id", req.UserID)
	return nil
}
