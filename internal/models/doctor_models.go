package models

type Doctor struct {
	Base
	UserID          uint   `json:"user_id"`
	Specialization  string `json:"specialization"`
	ExperienceYears int    `json:"experience_years"`
	Bio             string `json:"bio"`
}

type DoctorCreateRequest struct {
	UserID          uint   `json:"user_id"`
	Specialization  string `json:"specialization"`
	ExperienceYears int    `json:"experience_years"`
	Bio             string `json:"bio"`
}

type DoctorUpdateRequest struct {
	UserID          *uint   `json:"user_id"`
	Specialization  *string `json:"specialization"`
	ExperienceYears *int    `json:"experience_years"`
	Bio             *string `json:"bio"`
}
