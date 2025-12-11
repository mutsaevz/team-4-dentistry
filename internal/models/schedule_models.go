package models

import (
	"time"
)

type Schedule struct {
	Base
	DoctorID    uint      `json:"doctor_id" gorm:"not null;index"`
	Date        time.Time `json:"date" gorm:"not null;index"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	EndTime     time.Time `json:"end_time" gorm:"not null"`
	RoomNumber  int       `json:"room_number" gorm:"not null"`
	IsAvailable bool      `json:"is_available" gorm:"default:true"`
}

type ScheduleCreateRequest struct {
	DoctorID    uint      `json:"doctor_id" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	RoomNumber  int       `json:"room_number" validate:"required"`
	IsAvailable bool      `json:"is_available" validate:"omitempty"`
}

type ScheduleUpdateRequest struct {
	DoctorID    *uint      `json:"doctor_id,omitempty" validate:"omitempty"`
	Date        *time.Time `json:"date,omitempty" validate:"omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty" validate:"omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty" validate:"omitempty,gtfield=StartTime"`
	RoomNumber  *int       `json:"room_number,omitempty" validate:"omitempty"`
	IsAvailable *bool      `json:"is_available,omitempty" validate:"omitempty"`
}
