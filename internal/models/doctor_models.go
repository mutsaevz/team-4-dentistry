package models

type Doctor struct {
	Base
	UserID          uint   `json:"user_id" gorm:"not null;uniqueIndex"`
	Specialization  string `json:"specialization" gorm:"not null;size:100"`
	ExperienceYears int    `json:"experience_years" gorm:"not null;default:0"`
	Bio             string `json:"bio,omitempty" gorm:"type:text"`
}

type DoctorCreateRequest struct {
	UserID          uint   `json:"user_id" validate:"required"`
	Specialization  string `json:"specialization" validate:"required,max=100"`
	ExperienceYears int    `json:"experience_years" validate:"gte=0"`
	Bio             string `json:"bio,omitempty" validate:"max=2000"`
}

type DoctorUpdateRequest struct {
	UserID          *uint   `json:"user_id,omitempty" validate:"omitempty"`
	Specialization  *string `json:"specialization,omitempty" validate:"omitempty,max=100"`
	ExperienceYears *int    `json:"experience_years,omitempty" validate:"omitempty,gte=0"`
	Bio             *string `json:"bio,omitempty" validate:"omitempty,max=2000"`
}
