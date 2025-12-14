package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/config"
	"github.com/mutsaevz/team-4-dentistry/internal/loggers"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
	"github.com/mutsaevz/team-4-dentistry/internal/seed"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
	"github.com/mutsaevz/team-4-dentistry/internal/transports"
)

func main() {
	logger := loggers.InitLogger()

	db := config.SetUpDatabaseConnection(logger)

	userRepo := repository.NewUserRepository(db, logger)
	serviceRepo := repository.NewServiceRepository(db, logger)
	doctorRepo := repository.NewDoctorRepository(db, logger)
	scheduleRepo := repository.NewScheduleRepository(db, logger)
	reviewRepo := repository.NewReviewRepository(db, logger)
	patientRecordRepo := repository.NewPatientRecordRepo(db, logger)
	recommendationRepo := repository.NewRecommendationRepository(db, logger)

	if err := db.AutoMigrate(
		&models.Appointment{},
		&models.Doctor{},
		&models.PatientRecord{},
		&models.Recommendation{},
		&models.Review{},
		&models.Schedule{},
		&models.Service{},
		&models.User{},
	); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	if err := seed.SeedAdmin(userRepo, logger); err != nil {
		logger.Error("Не удалось заполнить административную панель", "error", err)
		os.Exit(1)
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	jwtCfg := services.JWTConfig{
		SecretKey:      secret,
		AccessTokenTTL: time.Hour * 24,
	}

	userService := services.NewUserService(userRepo, logger)
	servService := services.NewServService(serviceRepo, logger)
	doctorService := services.NewDoctorService(doctorRepo, serviceRepo, scheduleRepo, logger)
	authService := services.NewAuthService(userRepo, jwtCfg, logger)
	scheduleService := services.NewScheduleService(scheduleRepo, doctorRepo, logger)
	reviewService := services.NewReviewService(reviewRepo, doctorRepo, userRepo, logger)
	patientRecordService := services.NewPatientRecordService(patientRecordRepo, logger)
	recommendationService := services.NewRecommendationService(
		recommendationRepo,
		userRepo,
		serviceRepo,
		logger,
	)

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	transports.RegisterRoutes(
		r,
		logger,
		servService,
		userService,
		authService,
		jwtCfg,
		recommendationService,
		doctorService,
		scheduleService,
		reviewService,
		patientRecordService,
	)

	addr := ":8080"

	if err := r.Run(addr); err != nil {
		logger.Error("ошибка при запуске сервера", "addr", addr, "error", err)
		os.Exit(1)
	}
}
