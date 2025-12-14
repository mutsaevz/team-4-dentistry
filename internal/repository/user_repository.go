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
	db *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) UserRepository {
	return &gormUserRepository{db: db, logger: logger}
}

func (r *gormUserRepository) Create(user *models.User) error {
	if user == nil {
		return constants.User_IS_nil
	}

	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormUserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *gormUserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *gormUserRepository) List(offset, limit int) ([]models.User, error) {
	var users []models.User

	if err := r.db.
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *gormUserRepository) Update(user *models.User) error {
	if user == nil {
		return constants.User_IS_nil
	}

	if err := r.db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormUserRepository) Delete(id uint) error {

	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	return nil
}
