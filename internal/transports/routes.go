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
) {
	authHandler := NewAuthHandler(authService, userService)
	serviceHandler := NewServiceHandler(servService)
	userHandler := NewUserHandler(userService)

	authHandler.RegisterRoutes(router)

	api := router.Group("/api")

	serviceHandler.RegisterRoutes(api)
	userHandler.RegisterRoutes(api)
}
