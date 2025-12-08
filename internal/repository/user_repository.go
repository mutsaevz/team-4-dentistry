package repository

import (
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error

	GetByID(id uint) (*models.User, error)

	Update(user *models.User) error

	Delete(id uint) error
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *models.User) error {
	if user == nil {
		return nil
	}

	if err := r.db.Create(&user).Error; err != nil {
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

func (r *gormUserRepository) Update(user *models.User) error {
	if user == nil {
		return nil
	}

	if err := r.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormUserRepository) Delete(id uint) error {
	var user models.User

	if err := r.db.Delete(&user, id).Error; err != nil {
		return err
	}

	return nil
}