package services

import (
	"context"
	"errors"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
)

type DoctorService interface {
	CreateDoctor(context.Context, models.DoctorCreateRequest) (*models.Doctor, error)

	GetDoctorByID(context.Context, uint) (*models.Doctor, error)

	ListDoctors(context.Context, models.DoctorQueryParams) ([]models.Doctor, error)

	UpdateDoctor(context.Context, uint, models.DoctorUpdateRequest) (*models.Doctor, error)

	DeleteDoctor(context.Context, uint) error

	GetDoctorServices(context.Context, uint) ([]models.Service, error)
}

type doctorService struct {
	doctors repository.DoctorRepository
	service repository.ServiceRepository
}

func NewDoctorService(
	doctors repository.DoctorRepository,
	service repository.ServiceRepository,
) DoctorService {
	return &doctorService{doctors: doctors, service: service}
}

func (s *doctorService) CreateDoctor(ctx context.Context, req models.DoctorCreateRequest) (*models.Doctor, error) {

	if err := s.ValidateCreateDoctor(req); err != nil {
		return nil, err
	}

	doctor := &models.Doctor{
		UserID:          req.UserID,
		Specializations: append([]string{}, req.Specializations...),
		ExperienceYears: req.ExperienceYears,
		Bio:             req.Bio,
		AvgRating:       0,
		RoomNumber:      req.RoomNumber,
	}

	if err := s.doctors.Create(ctx, doctor); err != nil {
		return nil, err
	}

	return doctor, nil
}

func (s *doctorService) GetDoctorByID(ctx context.Context, id uint) (*models.Doctor, error) {
	return s.doctors.GetByID(id, ctx)
}

func (s *doctorService) ListDoctors(ctx context.Context, params models.DoctorQueryParams) ([]models.Doctor, error) {
	return s.doctors.GetAll(params, ctx)
}

func (s *doctorService) GetDoctorServices(ctx context.Context, doctorID uint) ([]models.Service, error) {
	return s.service.GetServicesByDoctorID(ctx, doctorID)
}

func (s *doctorService) UpdateDoctor(ctx context.Context, id uint, req models.DoctorUpdateRequest) (*models.Doctor, error) {
	doctor, err := s.doctors.GetByID(id, ctx)
	if err != nil {
		return nil, err
	}

	if len(*req.Specializations) == 0 {
		doctor.Specializations = append(doctor.Specializations, *req.Specializations...)
	}

	if *req.ExperienceYears != 0 {
		doctor.ExperienceYears = *req.ExperienceYears
	}

	if err := s.doctors.Update(ctx, doctor); err != nil {
		return nil, err
	}

	return doctor, nil
}

func (s *doctorService) DeleteDoctor(ctx context.Context, id uint) error {
	return s.doctors.Delete(ctx, id)
}

func (s *doctorService) ValidateCreateDoctor(req models.DoctorCreateRequest) error {
	if req.UserID <= 0 {
		return errors.New("")
	}
	if len(req.Specializations) == 0 {
		return errors.New("")
	}
	if req.RoomNumber <= 0 {
		return errors.New("")
	}
	if req.ExperienceYears < 0 {
		return errors.New("")
	}
	if req.Bio == "" {
		return errors.New("")
	}

	return nil
}
