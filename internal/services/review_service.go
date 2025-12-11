package services

import (
	"context"
	"errors"

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
}

func NewReviewService(review repository.ReviewRepository,
	doctor repository.DoctorRepository,
	patient repository.UserRepository) ReviewService {
	return &reviewService{
		review:  review,
		doctor:  doctor,
		patient: patient,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, req models.ReviewCreateRequest) (*models.Review, error) {

	if err := s.ValidateCreateReview(req); err != nil {
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
		return nil, err
	}

	return &review, nil
}

func (s *reviewService) GetByID(ctx context.Context, id uint) (*models.Review, error) {
	return s.review.GetByID(ctx, id)
}

func (s *reviewService) UpdateReview(ctx context.Context, id uint, req models.ReviewUpdateRequest) (*models.Review, error) {
	review, err := s.review.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Rating != nil {
		review.Rating = *req.Rating
	}

	if req.Comment != nil {
		review.Comment = *req.Comment
	}

	if err := s.review.Update(review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) DeleteReview(ctx context.Context, id uint) error {
	return s.review.Delete(ctx, id)
}

func (s *reviewService) GetDoctorReviews(ctx context.Context, doctorID uint) ([]models.Review, error) {
	return s.review.GetByDoctorID(ctx, doctorID)
}

func (s *reviewService) GetPatientReviews(ctx context.Context, patientID uint) ([]models.Review, error) {
	return s.review.GetByPatientID(ctx, patientID)
}

func (s *reviewService) ValidateCreateReview(req models.ReviewCreateRequest) error {
	if req.AppointmentID == 0 {
		return errors.New("")
	}
	if req.UserID == 0 {
		return errors.New("")
	}
	if req.DoctorID == 0 {
		return errors.New("")
	}
	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("")
	}
	return nil
}
