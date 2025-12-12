package transports

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	servService services.ServService,
	userService services.UserService,
	authService services.AuthService,
	jwtCfg services.JWTConfig,
	recService services.RecommendationService,
	docService services.DoctorService,
	scheduleService services.ScheduleService,
	reviewService services.ReviewService,
	patientRecordService services.PatientRecordService,
) {

	//patientRecordHandler := NewPatientRecordHandler(patientRecordService)

	api := router.Group("/api")

	// ---AUTH----
	authHandler := NewAuthHandler(authService, userService, logger)
	authHandler.RegisterRoutes(api, jwtCfg)

	protected := api.Group("")
	protected.Use(AuthMiddleware(jwtCfg))

	//----USER----
	userHandler := NewUserHandler(userService, logger)
	userHandler.RegisterRoutes(api)

	//----Service----
	serviceHandler := NewServiceHandler(servService, logger)
	serviceHandler.RegisterRoutes(api)

	//----RECOMMENDATION----
	recHandler := NewRecommendationHandler(recService, logger)
	recHandler.RegisterRoutes(api)

	//----SCHEDULE----
	scheduleHandler := NewScheduleHandler(scheduleService, logger)
	scheduleHandler.RegisterRoutes(api)

	//----DOCTOR----
	docHandler := NewDoctorHandler(docService, servService, scheduleService, reviewService, logger)
	docHandler.RegisterRoutes(api)

	//----REVIEW----
	reviewHandler := NewReviewHandler(reviewService, logger)
	reviewHandler.RegisterRoutes(api)
}
