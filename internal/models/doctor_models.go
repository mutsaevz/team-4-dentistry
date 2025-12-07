package models

type Doctor struct {
	Base
	UserID          uint     `json:"user_id" gorm:"not null;uniqueIndex"`
	Specializations []string `json:"specializations" gorm:"type:text[]"`
	ExperienceYears int      `json:"experience_years" gorm:"not null;default:0"`
	Bio             string   `json:"bio,omitempty" gorm:"type:text"`
	AvgRating       float64  `json:"avg_rating"`

	Shedules []Shedule `json:"-"`
	Reviews  []Review  `json:"-"`
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

type DoctorQueryParams struct {
	Specialization  string
	ExperienceYears int
	AvgRating       float64

	FilOr bool
}
