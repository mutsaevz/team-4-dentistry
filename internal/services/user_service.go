package services

import (
	"errors"
	"strings"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("пользователь не найден")

type UserService interface {
	CreateUser(req models.UserCreateRequest) (*models.User, error)

	GetUserById(id uint) (*models.User, error)

	UpdateUser(id uint, req models.UserUpdateRequest) (*models.User, error)

	DeleteUser(id uint) error
}

type userService struct {
	users repository.UserRepository
}

func NewUserService(
	users repository.UserRepository,
) UserService {
	return &userService{users: users}
}

func (s *userService) CreateUser(
	req models.UserCreateRequest,
) (*models.User, error) {

	if err := s.ValidateCreateUser(req); err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
		Email:     strings.TrimSpace(req.Email),
		Phone:     strings.TrimSpace(req.Phone),
		Password:  strings.TrimSpace(req.Password),
		Role:      req.Role,
	}

	if err := s.users.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserById(id uint) (*models.User, error) {
	user, err := s.users.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(
	id uint, req models.UserUpdateRequest,
) (*models.User, error) {
	user, err := s.users.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := s.ApplyUserUpdate(user, req); err != nil {
		return nil, err
	}

	if err := s.users.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	if _, err := s.users.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if err := s.users.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *userService) ValidateCreateUser(req models.UserCreateRequest) error {
	if req.FirstName == "" {
		return errors.New("имя не должно быть пустым")
	}

	if req.LastName == "" {
		return errors.New("фамилия не должна быть пустой")
	}

	if req.Password == "" {
		return errors.New("пароль не должен быть пустым")
	}

	if req.Phone == "" {
		return errors.New("телефон не должен быть пустым")
	}

	return nil
}

func (s *userService) ApplyUserUpdate(
	user *models.User,
	req models.UserUpdateRequest,
) error {
	if req.FirstName != nil {
		trimmed := strings.TrimSpace(*req.FirstName)

		if trimmed == "" {
			return errors.New("Имя не должно быть пустым")
		}

		user.FirstName = trimmed
	}

	if req.LastName != nil {
		trimmed := strings.TrimSpace(*req.LastName)
		if trimmed == "" {
			return errors.New("фамилия не должна быть пустой")
		}

		user.LastName = trimmed
	}

	if req.Password != nil {
		trimmed := strings.TrimSpace(*req.Password)
		if trimmed == "" {
			return errors.New("пароль не должен быть пустым")
		}

		user.Password = *req.Password
	}
	if req.Phone != nil {
		trimmed := strings.TrimSpace(*req.Phone)
		if trimmed == "" {
			return errors.New("телефон не должен быть пустым")
		}
		user.Phone = trimmed
	}
	return nil
}
