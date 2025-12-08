package models

import "time"

type Appointment struct {
	Base
	PatientID uint `json:"patient_id" gorm:"not null"`
	//Patient            *User      `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	DentistID uint    `json:"dentist_id" gorm:"not null"`
	Dentist   *Doctor `json:"dentist,omitempty" gorm:"foreignKey:DentistID"`
	ServiceID uint    `json:"service_id,omitempty"`
	//Service            *Service   `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	StartAt            time.Time  `json:"start_at" gorm:"not null"`
	EndAt              time.Time  `json:"end_at" gorm:"not null"`
	DurationMinutes    int        `json:"duration_minutes,omitempty"`
	Status             string     `json:"status" gorm:"type:varchar(50);default:'scheduled'"`
	Price              int        `json:"price_cents,omitempty"`
	Paid               bool       `json:"paid" gorm:"default:false"`
	Room               string     `json:"room,omitempty"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty"`
}

type AppointmentCreateRequest struct {
	PatientID       uint      `json:"patient_id" validate:"required"`
	DentistID       uint      `json:"dentist_id" validate:"required"`
	ServiceID       uint     `json:"service_id,omitempty" validate:"omitempty"`
	StartAt         time.Time `json:"start_at" validate:"required"`
	EndAt           time.Time `json:"end_at" validate:"required,gtfield=StartAt"`
	DurationMinutes int       `json:"duration_minutes,omitempty" validate:"omitempty,min=1"`
	Price           int       `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Room            string    `json:"room,omitempty" validate:"omitempty"`
	Notes           string    `json:"notes,omitempty" validate:"omitempty"`
}

type AppointmentUpdateRequest struct {
	PatientID       *uint      `json:"patient_id,omitempty" validate:"omitempty"`
	DentistID       *uint      `json:"dentist_id,omitempty" validate:"omitempty"`
	ServiceID       *uint      `json:"service_id,omitempty" validate:"omitempty"`
	StartAt         *time.Time `json:"start_at,omitempty" validate:"omitempty"`
	EndAt           *time.Time `json:"end_at,omitempty" validate:"omitempty,gtfield=StartAt"`
	DurationMinutes *int       `json:"duration_minutes,omitempty" validate:"omitempty,min=1"`
	Price           *int       `json:"price_cents,omitempty" validate:"omitempty,min=0"`
	Room            *string    `json:"room,omitempty" validate:"omitempty"`
	Notes           *string    `json:"notes,omitempty" validate:"omitempty"`
}
