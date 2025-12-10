package models

import "time"

type Role string

const (
	Admin   Role = "admin"
	Doc     Role = "doctor"
	Patient Role = "patient"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type User struct {
	Base
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Password      string    `json:"-"`
	Role          Role      `json:"role"`
	Gender        Gender    `json:"gender"`
	EmailVerified bool      `json:"email_verified"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	IsActive      bool      `json:"is_active"`
}

type UserCreateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	Role        Role   `json:"role"`
	Gender      Gender `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

type UserUpdateRequest struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Password    *string `json:"password"`
	Role        *Role   `json:"role"`
	Gender      *Gender `json:"gender"`
	DateOfBirth *string `json:"date_of_birth"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
