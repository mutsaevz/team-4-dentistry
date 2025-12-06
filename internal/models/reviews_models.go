package models

type Review struct {
	Base
	AppointmentID uint   `json:"appointment_id"`
	UserID        uint   `json:"user_id"`
	DoctorID      uint   `json:"doctor_id"`
	Rating        int    `json:"rating"`
	Comment       string `json:"comment"`
}

type ReviewCreateRequest struct {
	AppointmentID uint   `json:"appointment_id"`
	UserID        uint   `json:"user_id"`
	DoctorID      uint   `json:"doctor_id"`
	Rating        int    `json:"rating"`
	Comment       string `json:"comment"`
}

type ReviewUpdateRequest struct {
	AppointmentID *uint   `json:"appointment_id"`
	UserID        *uint   `json:"user_id"`
	DoctorID      *uint   `json:"doctor_id"`
	Rating        *int    `json:"rating"`
	Comment       *string `json:"comment"`
}
