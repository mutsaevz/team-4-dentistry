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
) {
	serviceHandler := NewServiceHandler(servService)
	userHandler := NewUserHandler(userService)

	serviceHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)
}
