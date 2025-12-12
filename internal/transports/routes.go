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
	api := router.Group("/api")

	authHandler := NewAuthHandler(authService, userService)
	authHandler.RegisterRoutes(api, jwtCfg)

	// Группа где jwt обязателен
	protected := api.Group("")
	protected.Use(AuthMiddleware(jwtCfg))

	// Наши публичные услуги
	serviceHandler := NewServiceHandler(servService)
	servicePublic := api.Group("/services")
	servicePublic.GET("", serviceHandler.List)
	servicePublic.GET("/:id", serviceHandler.GetByID)

	// Защищенные
	serviceAdmin := protected.Group("/services")
	serviceAdmin.Use(RequireRole("admin"))
	serviceAdmin.POST("", serviceHandler.Create)
	serviceAdmin.PUT("/:id", serviceHandler.Update)
	serviceAdmin.DELETE("/:id", serviceHandler.Delete)

	// публичные

	docHandler := NewDoctorHandler(docService, servService, scheduleService, reviewService)
	docPublic := api.Group("/doctors")
	docPublic.GET("", docHandler.ListDoctors)
	docPublic.GET("/:id", docHandler.GetDoctorByID)
	docPublic.GET("/:id/reviews", docHandler.GetDoctorReviews)
	docPublic.GET("/:id/services", docHandler.ListDoctorServices)
	docPublic.GET("/:id/schedules/available", docHandler.GetAvailableSlots)

	// защищенные

	docAdmin := protected.Group("/doctors")
	docAdmin.Use(RequireRole("admin"))
	docAdmin.POST("", docHandler.CreateDoctor)
	docAdmin.PATCH("/:id", docHandler.UpdateDoctor)
	docAdmin.DELETE("/:id", docHandler.DeleteDoctor)
	docAdmin.GET("/:id/schedules", docHandler.ListSchedules)

	// все остальное защищенное, можем потом изменить по желанию

	// Users только admin
	userHandler := NewUserHandler(userService)
	userHandler.RegisterRoutes(protected)

	// Schedules только admin
	scheduleHandler := NewScheduleHandler(scheduleService)
	scheduleHandler.RegisterRoutes(protected)

	// Recommendations пациент читает "my", доктор/админ создаёт/удаляет
	recHandler := NewRecommendationHandler(recService)
	recs := protected.Group("/recommendations")
	recs.POST("", recHandler.Create)
	recs.GET("/my", recHandler.ListMy)
	recs.DELETE("/:id", recHandler.Delete)

	// Patient records как минимум админ/доктор
	patientRecordHandler := NewPatientRecordHandler(patientRecordService)
	records := protected.Group("/patient-records")
	records.Use(RequireRole("admin", "doctor"))
	records.POST("", patientRecordHandler.Create)
	records.GET("", patientRecordHandler.GetAll)
	records.GET("/:id", patientRecordHandler.GetByID)
	records.PATCH("/:id", patientRecordHandler.Update)
	records.DELETE("/:id", patientRecordHandler.Delete)

	// При код ревью в
	// Reviews  у нас в handler обнаружились баги с Param("id") vs doctor_id/patient_id,
	// поэтому лучше пока пользоваться публичным /api/doctors/:id/reviews.
	// Исправим баг потом
	// reviewHandler := NewReviewHandler(reviewService)
	// reviewHandler.RegisterRoutes(protected)
}
