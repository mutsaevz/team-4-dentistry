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

	//patientRecordHandler := NewPatientRecordHandler(patientRecordService)

	api := router.Group("/api")

	// ---AUTH----
	authHandler := NewAuthHandler(authService, userService)
	authHandler.RegisterRoutes(api, jwtCfg)

	protected := api.Group("")
	protected.Use(AuthMiddleware(jwtCfg))

	//----USER----
	userHandler := NewUserHandler(userService)
	userHandler.RegisterRoutes(api)

	//----Service----
	serviceHandler := NewServiceHandler(servService)
	serviceHandler.RegisterRoutes(api)

	//----RECOMMENDATION----
	recHandler := NewRecommendationHandler(recService)
	recHandler.RegisterRoutes(api)

	//----SCHEDULE----
	scheduleHandler := NewScheduleHandler(scheduleService)
	scheduleHandler.RegisterRoutes(api)

	//----DOCTOR----
	docHandler := NewDoctorHandler(docService, servService, scheduleService, reviewService)
	docHandler.RegisterRoutes(api)

	//----REVIEW----
	reviewHandler := NewReviewHandler(reviewService)
	reviewHandler.RegisterRoutes(api)
}
