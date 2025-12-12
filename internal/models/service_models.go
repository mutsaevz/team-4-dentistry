package models

type Service struct {
	Base
	Name        string  `json:"name"`
	DoctorID    uint    `json:"doctor_id" gorm:"not null;index"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Duration    int     `json:"duration"`
	Price       float64 `json:"price"`
}

type ServiceCreateRequest struct {
	DoctorID    uint    `json:"doctor_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Duration    int     `json:"duration"`
	Price       float64 `json:"price"`
}

type ServiceUpdateRequest struct {
	DoctorID    *uint    `json:"doctor_id,omitempty"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Category    *string  `json:"category"`
	Duration    *int     `json:"duration"`
	Price       *float64 `json:"price"`
}
