package transports

import (
	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

func RegisterRoutes(
	router *gin.Engine,
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
	authHandler := NewAuthHandler(authService, userService)
	serviceHandler := NewServiceHandler(servService)
	userHandler := NewUserHandler(userService)
	recHandler := NewRecommendationHandler(recService)
	docHandler := NewDoctorHandler(docService, servService, scheduleService, reviewService)
	scheduleHandler := NewScheduleHandler(scheduleService)
	reviewHandler := NewReviewHandler(reviewService)
	//patientRecordHandler := NewPatientRecordHandler(patientRecordService)

	authHandler.RegisterRoutes(router)
	recHandler.RegisterRoutes(router)
	scheduleHandler.RegisterRoutes(router)

	api := router.Group("/api")
	api.Use(AuthMiddleware(jwtCfg))

	serviceHandler.RegisterRoutes(api)
	userHandler.RegisterRoutes(api)
	docHandler.RegisterRoutes(api)
	reviewHandler.RegisterRoutes(api)
}
