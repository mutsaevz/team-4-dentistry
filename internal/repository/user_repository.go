package repository

import (
	"log/slog"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error

	GetByID(id uint) (*models.User, error)

	GetByEmail(email string) (*models.User, error)

	List(offset, limit int) ([]models.User, error)

	Update(user *models.User) error

	Delete(id uint) error
}

type gormUserRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) UserRepository {
	return &gormUserRepository{db: db, logger: logger}
}

func (r *gormUserRepository) Create(user *models.User) error {
	if user == nil {
		r.logger.Warn("попытка создать nil user")
		return constants.User_IS_nil
	}

	r.logger.Debug("создание user в репозитории", "email", user.Email)

	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error("ошибка при создании user", "error", err, "email", user.Email)
		return err
	}

	r.logger.Info("user создан", "user_id", user.ID, "email", user.Email)
	return nil
}

func (r *gormUserRepository) GetByID(id uint) (*models.User, error) {
	r.logger.Debug("получение user по ID", "user_id", id)
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil {
		r.logger.Error("ошибка при получении user по ID", "error", err, "user_id", id)
		return nil, err
	}

	r.logger.Info("user получен по ID", "user_id", id)
	return &user, nil
}

func (r *gormUserRepository) GetByEmail(email string) (*models.User, error) {
	r.logger.Debug("получение user по email", "email", email)
	var user models.User

	if err := r.db.
		Where("email = ?", email).
		First(&user).Error; err != nil {
		r.logger.Error("ошибка при получении user по email", "error", err, "email", email)
		return nil, err
	}

	r.logger.Info("user получен по email", "user_id", user.ID, "email", email)
	return &user, nil
}

func (r *gormUserRepository) List(offset, limit int) ([]models.User, error) {
	r.logger.Debug("получение списка users", "offset", offset, "limit", limit)
	var users []models.User

	if err := r.db.
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		r.logger.Error("ошибка при получении списка users", "error", err, "offset", offset, "limit", limit)
		return nil, err
	}

	r.logger.Info("список users получен", "count", len(users))
	return users, nil
}

func (r *gormUserRepository) Update(user *models.User) error {
	if user == nil {
		r.logger.Warn("попытка обновить nil user")
		return constants.User_IS_nil
	}

	r.logger.Debug("обновление user", "user_id", user.ID)

	if err := r.db.Save(user).Error; err != nil {
		r.logger.Error("ошибка при обновлении user", "error", err, "user_id", user.ID)
		return err
	}

	r.logger.Info("user успешно обновлен", "user_id", user.ID)
	return nil
}

func (r *gormUserRepository) Delete(id uint) error {
	r.logger.Debug("удаление user по ID", "user_id", id)

	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		r.logger.Error("ошибка при удалении user", "error", err, "user_id", id)
		return err
	}

	r.logger.Info("user успешно удален", "user_id", id)
	return nil
}
