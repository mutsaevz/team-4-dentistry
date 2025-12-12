package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/config"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
	"github.com/mutsaevz/team-4-dentistry/internal/seed"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
	"github.com/mutsaevz/team-4-dentistry/internal/transports"
)

func main() {

	db := config.SetUpDatabaseConnection()

	userRepo := repository.NewUserRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	// doctorRepo := repository.NewDoctorRepository(db)
	// scheduleRepo := repository.NewScheduleRepository(db)
	// reviewRepo := repository.NewReviewRepository(db)
	// patientRecordRepo := repository.NewPatientRecordRepo(db)
	recommendationRepo := repository.NewRecommendationRepository(db)

	if err := seed.SeedAdmin(userRepo); err != nil {
		log.Fatalf("Не удалось заполнить административную панель: %v", err)
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	jwtCfg := services.JWTConfig{
		SecretKey:      secret,
		AccessTokenTTL: time.Hour * 24,
	}

	userService := services.NewUserService(userRepo)
	servService := services.NewServService(serviceRepo)
	//doctorService := services.NewDoctorService(doctorRepo, serviceRepo)
	authService := services.NewAuthService(userRepo, jwtCfg)
	//scheduleService := services.NewScheduleService(scheduleRepo, doctorRepo)
	//reviewService := services.NewReviewService(reviewRepo, doctorRepo, userRepo)
	//patientRecordService := services.NewPatientRecordService(patientRecordRepo)
	recommendationService := services.NewRecommendationService(
		recommendationRepo,
		userRepo,
		serviceRepo,
	)

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	transports.RegisterRoutes(
		r,
		servService,
		userService,
		authService,
		jwtCfg,
		recommendationService,
	)

	addr := ":8080"

	if err := r.Run(addr); err != nil {
		log.Fatalf("ошибка при запуске сервера %s: %v", addr, err)
	}
}
