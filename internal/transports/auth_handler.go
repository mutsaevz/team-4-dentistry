package transports

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/services"
)

type AuthHandler struct {
	auth  services.AuthService
	users services.UserService
	logger *slog.Logger
}

func NewAuthHandler(auth services.AuthService, users services.UserService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		auth:  auth,
		users: users,
		logger: logger,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup, jwtCfg services.JWTConfig) {
	auth := r.Group("/auth")
	// ----публичные----
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	// -----нужна-авторизация----
	protected := auth.Group("")
	protected.Use(AuthMiddleware(jwtCfg))
	protected.POST("/refresh", h.Refresh)
	protected.POST("/logout", h.Logout)
	protected.GET("/me", h.Me)
	protected.PUT("/me", h.UpdateMe)
	protected.PUT("/me/password", h.ChangePassword)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	req.Role = models.Patient

	user, err := h.users.CreateUser(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	token, err := h.auth.Login(req.Email, req.Password)

	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "неверный email или пароль",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}

func (h *AuthHandler) Me(c *gin.Context) {
	idVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	userID, ok := idVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный userID в context",
		})
		return
	}
	user, err := h.users.GetUserById(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	idVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	userID, ok := idVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный userID в context",
		})
		return
	}
	var req models.UserUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	user, err := h.users.UpdateUser(userID, req)

	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	idVal, exists := c.Get("userID")
	roleVal, existsRole := c.Get("userRole")

	if !exists || !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	userdID, okID := idVal.(uint)
	role, okRole := roleVal.(string)

	if !okID || !okRole {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный данные",
		})
		return
	}

	token, err := h.auth.GenerateToken(userdID, role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	idVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
		return
	}

	userID, ok := idVal.(uint)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "некорректный userID в context",
		})
		return
	}

	var req models.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	if err := h.users.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
