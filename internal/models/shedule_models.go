package models

import (
	"time"
)

type Shedule struct {
	Base
	DoctorID   uint      `json:"doctor_id" gorm:"not null;index"`
	Date       time.Time `json:"date" gorm:"not null;index"`
	StartTime  time.Time `json:"start_time" gorm:"not null"`
	EndTime    time.Time `json:"end_time" gorm:"not null"`
	RoomNumber string    `json:"room_number" gorm:"not null"`
}

type SheduleCreateRequest struct {
	DoctorID   uint      `json:"doctor_id" validate:"required"`
	Date       time.Time `json:"date" validate:"required"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	RoomNumber string    `json:"room_number" validate:"required"`
}

type SheduleUpdateRequest struct {
	DoctorID   *uint      `json:"doctor_id,omitempty" validate:"omitempty"`
	Date       *time.Time `json:"date,omitempty" validate:"omitempty"`
	StartTime  *time.Time `json:"start_time,omitempty" validate:"omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty" validate:"omitempty,gtfield=StartTime"`
	RoomNumber *string    `json:"room_number,omitempty" validate:"omitempty"`
}
