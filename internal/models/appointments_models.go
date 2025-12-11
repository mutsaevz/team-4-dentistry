package models

import "time"

type Appointment struct {
	Base
	PatientID   uint       `json:"patient_id" gorm:"not null"`
	Patient     *User      `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	DoctorID    uint       `json:"doctor_id" gorm:"not null"`
	Doctor      *Doctor    `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	ServiceID   uint       `json:"service_id,omitempty"`
	Service     *Service   `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	StartAt     time.Time  `json:"start_at" gorm:"not null"`
	EndAt       time.Time  `json:"end_at" gorm:"not null"`
	Status      string     `json:"status" gorm:"type:varchar(50);default:'scheduled'"`
	Price       float64    `json:"price_cents,omitempty"`
	Paid        bool       `json:"paid" gorm:"default:false"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
}

type AppointmentCreateRequest struct {
	PatientID uint      `json:"patient_id" validate:"required"`
	DoctorID  uint      `json:"doctor_id" validate:"required"`
	ServiceID uint      `json:"service_id,omitempty" validate:"omitempty"`
	StartAt   time.Time `json:"start_at" validate:"required"`
	Price     float64   `json:"price_cents,omitempty" validate:"omitempty,min=0.0"`
}

type AppointmentUpdateRequest struct {
	PatientID *uint      `json:"patient_id,omitempty" validate:"omitempty"`
	DoctorID  *uint      `json:"doctor_id,omitempty" validate:"omitempty"`
	ServiceID *uint      `json:"service_id,omitempty" validate:"omitempty"`
	StartAt   *time.Time `json:"start_at,omitempty" validate:"omitempty"`
	Price     *float64   `json:"price_cents,omitempty" validate:"omitempty,min=0.0"`
}
