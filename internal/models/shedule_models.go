package models

import (
	"time"
)

type Shedule struct {
	Base
	DoctorID   uint      `json:"doctor_id"`
	Date       time.Time `json:"date"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	RoomNumber string    `json:"room_number"`
}

type SheduleCreateRequest struct {
	DoctorID   uint      `json:"doctor_id"`
	Date       time.Time `json:"date"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	RoomNumber string    `json:"room_number"`
}

type SheduleUpdateRequest struct {
	DoctorID   *uint      `json:"doctor_id"`
	Date       *time.Time `json:"date"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
	RoomNumber *string    `json:"room_number"`
}
