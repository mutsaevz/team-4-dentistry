package models

type Review struct {
	Base
	AppointmentID uint   `json:"appointment_id" gorm:"not null;index"`
	UserID        uint   `json:"user_id" gorm:"not null;index"`
	DoctorID      uint   `json:"doctor_id" gorm:"not null;index"`
	Rating        int    `json:"rating" gorm:"not null;check:rating >= 0 AND rating <= 5"`
	Comment       string `json:"comment,omitempty" gorm:"type:text"`
}

type ReviewCreateRequest struct {
	AppointmentID uint   `json:"appointment_id" validate:"required"`
	UserID        uint   `json:"user_id" validate:"required"`
	DoctorID      uint   `json:"doctor_id" validate:"required"`
	Rating        int    `json:"rating" validate:"required,gte=0,lte=5"`
	Comment       string `json:"comment,omitempty" validate:"max=2000"`
}

type ReviewUpdateRequest struct {
	AppointmentID *uint   `json:"appointment_id,omitempty" validate:"omitempty"`
	UserID        *uint   `json:"user_id,omitempty" validate:"omitempty"`
	DoctorID      *uint   `json:"doctor_id,omitempty" validate:"omitempty"`
	Rating        *int    `json:"rating,omitempty" validate:"omitempty,gte=0,lte=5"`
	Comment       *string `json:"comment,omitempty" validate:"omitempty,max=2000"`
}
