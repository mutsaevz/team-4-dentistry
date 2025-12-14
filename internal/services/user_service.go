package services

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("пользователь не найден")

type UserService interface {
	CreateUser(req models.UserCreateRequest) (*models.User, error)

	GetUserById(id uint) (*models.User, error)

	ListUsers(offset, limit int) ([]models.User, error)

	UpdateUser(id uint, req models.UserUpdateRequest) (*models.User, error)

	DeleteUser(id uint) error

	ChangePassword(userID uint, oldPassword, newPassword string) error
}

type userService struct {
	users  repository.UserRepository
	logger *slog.Logger
}

func NewUserService(
	users repository.UserRepository,
	logger *slog.Logger,
) UserService {
	return &userService{users: users, logger: logger}
}

func hashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 14)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func checkPassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func (s *userService) CreateUser(
	req models.UserCreateRequest,
) (*models.User, error) {
	s.logger.Debug("CreateUser called", "email", req.Email)

	if err := s.ValidateCreateUser(req); err != nil {
		s.logger.Error("validation failed for CreateUser", "error", err, "email", req.Email)
		return nil, err
	}

	hashed, err := hashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err, "email", req.Email)
		return nil, err
	}

	user := &models.User{
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
		Email:     strings.TrimSpace(req.Email),
		Phone:     strings.TrimSpace(req.Phone),
		Password:  hashed,
		Role:      models.Role(strings.TrimSpace(string(req.Role))),
	}

	if err := s.users.Create(user); err != nil {
		s.logger.Error("failed to create user in repo", "error", err, "email", req.Email)
		return nil, err
	}

	s.logger.Info("user created", "user_id", user.ID, "email", user.Email)
	return user, nil
}

func (s *userService) GetUserById(id uint) (*models.User, error) {
	s.logger.Debug("GetUserById called", "user_id", id)
	user, err := s.users.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found", "user_id", id)
			return nil, ErrUserNotFound
		}
		s.logger.Error("error getting user by id", "error", err, "user_id", id)
		return nil, err
	}

	s.logger.Info("user retrieved", "user_id", id)
	return user, nil
}

func (s *userService) ListUsers(offset, limit int) ([]models.User, error) {
	s.logger.Debug("ListUsers called", "offset", offset, "limit", limit)
	users, err := s.users.List(offset, limit)
	if err != nil {
		s.logger.Error("error listing users", "error", err)
		return nil, err
	}
	s.logger.Info("users listed", "count", len(users))
	return users, nil
}

func (s *userService) UpdateUser(
	id uint, req models.UserUpdateRequest,
) (*models.User, error) {
	s.logger.Debug("UpdateUser called", "user_id", id)
	user, err := s.users.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found for update", "user_id", id)
			return nil, ErrUserNotFound
		}
		s.logger.Error("error fetching user for update", "error", err, "user_id", id)
		return nil, err
	}

	if err := s.ApplyUserUpdate(user, req); err != nil {
		s.logger.Error("validation failed when applying user update", "error", err, "user_id", id)
		return nil, err
	}

	if err := s.users.Update(user); err != nil {
		s.logger.Error("failed to update user in repo", "error", err, "user_id", id)
		return nil, err
	}

	s.logger.Info("user updated", "user_id", id)
	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	s.logger.Debug("DeleteUser called", "user_id", id)
	if _, err := s.users.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found for delete", "user_id", id)
			return ErrUserNotFound
		}
		s.logger.Error("error fetching user for delete", "error", err, "user_id", id)
		return err
	}

	if err := s.users.Delete(id); err != nil {
		s.logger.Error("failed to delete user in repo", "error", err, "user_id", id)
		return err
	}

	s.logger.Info("user deleted", "user_id", id)
	return nil
}

func (s *userService) ChangePassword(
	userID uint,
	oldPassword,
	newPassword string,
) error {
	s.logger.Debug("ChangePassword called", "user_id", userID)
	user, err := s.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found for ChangePassword", "user_id", userID)
			return ErrUserNotFound
		}
		s.logger.Error("error fetching user for ChangePassword", "error", err, "user_id", userID)
		return err
	}

	if err := checkPassword(user.Password, oldPassword); err != nil {
		s.logger.Warn("old password mismatch", "user_id", userID)
		return errors.New("старый пароль неверен")
	}

	if strings.TrimSpace(newPassword) == "" {
		s.logger.Warn("new password is empty", "user_id", userID)
		return errors.New("новый пароль не должен быть пустым")
	}

	hashed, err := hashPassword(newPassword)
	if err != nil {
		s.logger.Error("failed to hash new password", "error", err, "user_id", userID)
		return err
	}

	user.Password = hashed

	if err := s.users.Update(user); err != nil {
		s.logger.Error("failed to update user password in repo", "error", err, "user_id", userID)
		return err
	}

	s.logger.Info("password changed", "user_id", userID)
	return nil
}

func (s *userService) ValidateCreateUser(req models.UserCreateRequest) error {
	if req.FirstName == "" {
		return errors.New("имя не должно быть пустым")
	}

	if req.LastName == "" {
		return errors.New("фамилия не должна быть пустой")
	}

	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email не должен быть пустым")
	}

	if req.Password == "" {
		return errors.New("пароль не должен быть пустым")
	}

	if req.Phone == "" {
		return errors.New("телефон не должен быть пустым")
	}

	if strings.TrimSpace(string(req.Role)) == "" {
		return errors.New("роль не должна быть пустой")
	}

	role := strings.TrimSpace(string(req.Role))

	switch role {
	case "admin", "doctor", "patient":
	default:
		return errors.New("некорректная роль")
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

		hashed, err := hashPassword(trimmed)
		if err != nil {
			return err
		}
		user.Password = hashed
	}
	if req.Phone != nil {
		trimmed := strings.TrimSpace(*req.Phone)
		if trimmed == "" {
			return errors.New("телефон не должен быть пустым")
		}
		user.Phone = trimmed
	}

	if req.Email != nil {
		trimmed := strings.TrimSpace(*req.Email)
		if trimmed == "" {
			return errors.New("email не должен быть пустым")
		}
		user.Email = trimmed
	}

	if req.Role != nil {
		trimmed := strings.TrimSpace(string(*req.Role))
		if trimmed == "" {
			return errors.New("Role не должен быть пустым")
		}
		switch trimmed {
		case "admin", "doctor", "patient":
			user.Role = models.Role(trimmed)
		default:
			return errors.New("некорректная роль")
		}

	}

	return nil
}
