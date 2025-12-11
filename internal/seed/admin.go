package seed

import (
	"log"
	"os"

	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(userRepo repository.UserRepository) error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		log.Println("[seed admin] ADMIN_EMAIL или ADMIN_PASSWORD не установлены")
		return nil
	}

	existing, _ := userRepo.GetByEmail(adminEmail)

	if existing != nil {
		log.Println("[seed admin] админ уже существует — skip")
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &models.User{
		Email:    adminEmail,
		Password: string(hash),
		Role:     "admin",
	}

	if err := userRepo.Create(admin); err != nil {
		return err
	}

	log.Printf("[seed admin] админ создан: %s\n", adminEmail)

	return nil
}
