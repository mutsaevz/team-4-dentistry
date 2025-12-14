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
	appointmentService services.AppointmentService,
) {
	api := router.Group("/api")

	// ---AUTH----
	authHandler := NewAuthHandler(authService, userService, logger)
	authHandler.RegisterRoutes(api, jwtCfg)

	// Группа где jwt обязателен
	protected := api.Group("")
	protected.Use(AuthMiddleware(jwtCfg))

	// Наши публичные услуги
	serviceHandler := NewServiceHandler(servService, logger)
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

	docHandler := NewDoctorHandler(docService, servService, scheduleService, reviewService, logger)
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

	// Users
	userHandler := NewUserHandler(userService, logger)
	userPublic := api.Group("/users")
	userPublic.POST("", userHandler.Create)

	// защищенные
	userAdmin := protected.Group("/users")
	userAdmin.Use(RequireRole("admin"))
	userAdmin.GET("", userHandler.List)
	userAdmin.GET("/:id", userHandler.GetByID)
	userAdmin.PUT("/:id", userHandler.Update)
	userAdmin.DELETE("/:id", userHandler.Delete)

	// Schedules только admin
	scheduleHandler := NewScheduleHandler(scheduleService, logger)
	scheduleHandler.RegisterRoutes(protected)

	// Recommendations пациент читает "my", доктор/админ создаёт/удаляет
	recHandler := NewRecommendationHandler(recService, logger)
	recs := protected.Group("/recommendations")
	recs.POST("", recHandler.Create)
	recs.GET("/my", recHandler.ListMy)
	recs.DELETE("/:id", recHandler.Delete)

	// Patient records как минимум админ/доктор
	patientRecordHandler := NewPatientRecordHandler(patientRecordService, logger)
	records := protected.Group("/patient-records")
	records.Use(RequireRole("admin", "doctor"))
	records.POST("", patientRecordHandler.Create)
	records.GET("", patientRecordHandler.GetAll)
	records.GET("/:id", patientRecordHandler.GetByID)
	records.PATCH("/:id", patientRecordHandler.Update)
	records.DELETE("/:id", patientRecordHandler.Delete)

	// Review
	reviewHandler := NewReviewHandler(reviewService, logger)
	reviewPublic := api.Group("/reviews")
	reviewPublic.POST("", reviewHandler.CreateReview)
	reviewPublic.GET("/doctor/:id", reviewHandler.GetDoctorReviews)
	reviewPublic.PUT("/:id", reviewHandler.UpdateReview)
	reviewPublic.DELETE("/:id", reviewHandler.DeleteReview)

	// Защищенные review
	reviewAdmin := protected.Group("/reviews")
	reviewAdmin.Use(RequireRole("admin"))
	reviewAdmin.GET("/:id", reviewHandler.GetReviewByID)
	reviewAdmin.GET("/patient/:patient_id", reviewHandler.GetPatientReviews)

	// Appointment
	appointmentHandler := NewAppointmentsHandler(appointmentService, logger)
	apPublic := api.Group("/appointments")
	apPublic.POST("", appointmentHandler.Create)
	apPublic.GET("/:id", appointmentHandler.GetByID)
	apPublic.PATCH("/:id", appointmentHandler.Update)
	apPublic.DELETE("/:id", appointmentHandler.Delete)
	apPublic.GET("/patients/:id", appointmentHandler.GetByPatientID)

	// защищенные
	apAdmin := protected.Group("/appointments")
	apAdmin.Use(RequireRole("admin"))
	apAdmin.GET("", appointmentHandler.GetAll)
}
